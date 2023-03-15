package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	net_http "net/http"
	"time"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
)

type checkWatcher struct {
	ctx    context.Context
	anchor *httpChecker
	probe  scout.ScoutDelegate
}

func newCheckWatcher(anchor *httpChecker, probe scout.ScoutDelegate) (*checkWatcher, error) {
	return &checkWatcher{
		ctx:    context.Background(),
		anchor: anchor,
		probe:  probe,
	}, nil
}

func (w *checkWatcher) Name() string {
	return w.probe.Name()
}

func (w *checkWatcher) Data() interface{} {
	return w.probe.Data()
}

func (w *checkWatcher) Helper() worker_queue.WatcherFunc {
	return func(i interface{}) (interface{}, error) {
		input := i.(*types.Endpoint)
		url := fmt.Sprintf("http://%s%s", input.ToString(), w.anchor.config.URL)
		req, err := net_http.NewRequest(w.anchor.config.Method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("creating the request for the health check failed: %w", err)
		}

		ctx, cancel := context.WithTimeout(w.ctx, time.Duration(w.anchor.config.RequestTimeout)*time.Millisecond)
		defer cancel()

		req.Header.Set("Connection", "close")
		req = req.WithContext(ctx)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("making the request for the health check failed: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode >= http.StatusInternalServerError {
			return nil, errors.New("remote service is not available at the moment")
		}

		return true, nil
	}
}

func (w *checkWatcher) Notify(a interface{}, b interface{}) {
	h, ok := a.(*types.Endpoint)
	if !ok {
		return
	}

	s, ok := b.(bool)
	if !ok {
		return
	}

	w.probe.Notify(h, s)
}

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

const (
	modelName             = "http"
	defaultRequestTimeout = 2000
)

type httpWatcher struct {
	config *HttpCheckerConfig
	ctx    context.Context
	item   scout.ScoutDelegate
}

func newHttpWatcher(conf *scout.WatcherConfig) (*httpWatcher, error) {
	var config *HttpCheckerConfig = nil
	config, ok := conf.Data.(*HttpCheckerConfig)
	if !ok {
		return nil, fmt.Errorf("http watcher invalid config: %v", conf)
	}

	if config.RequestTimeout == 0 {
		config.RequestTimeout = defaultRequestTimeout
	}

	return &httpWatcher{
		ctx:    context.Background(),
		config: config,
		item:   conf.Item,
	}, nil
}

func (w *httpWatcher) Name() string {
	return w.item.Name()
}

func (w *httpWatcher) Data() interface{} {
	return w.item.Data()
}

func (w *httpWatcher) Helper() worker_queue.WatcherFunc {
	return func(i interface{}) (interface{}, error) {
		input := i.(*types.Endpoint)
		url := fmt.Sprintf("http://%s%s", input.ToString(), w.config.URL)
		req, err := net_http.NewRequest(w.config.Method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("creating the request for the health check failed: %w", err)
		}

		ctx, cancel := context.WithTimeout(w.ctx, time.Duration(w.config.RequestTimeout)*time.Millisecond)
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

func (w *httpWatcher) Notify(a interface{}, b interface{}) {
	h, ok := a.(*types.Endpoint)
	if !ok {
		return
	}

	s, ok := b.(bool)
	if !ok {
		return
	}

	w.item.Notify(h, s)
}

func init() {
	scout.Register(modelName, func(conf *scout.WatcherConfig) (scout.ScoutWatcher, error) {
		return newHttpWatcher(conf)
	})
}

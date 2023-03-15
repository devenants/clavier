package tcp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
)

type checkWatcher struct {
	ctx    context.Context
	anchor *tcpChecker
	probe  scout.ScoutDelegate
}

func newCheckWatcher(anchor *tcpChecker, probe scout.ScoutDelegate) (*checkWatcher, error) {
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
		d := net.Dialer{Timeout: time.Duration(w.anchor.config.ConnectTimeout) * time.Millisecond}
		conn, err := d.Dial("tcp", input.ToString())
		if err != nil {
			return nil, fmt.Errorf("could not connect to server: %v", err)
		}

		defer conn.Close()

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

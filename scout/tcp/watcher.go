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

const (
	modelName             = "tcp"
	defaultConnectTimeout = 2000
)

type tcpWatcher struct {
	config *TcpCheckerConfig

	ctx  context.Context
	item scout.ScoutDelegate
}

func newTcpWatcher(conf *scout.WatcherConfig) (*tcpWatcher, error) {
	var config *TcpCheckerConfig = nil
	config, ok := conf.Data.(*TcpCheckerConfig)
	if !ok {
		return nil, fmt.Errorf("tcp watcher invalid config: %v", conf)
	}

	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = defaultConnectTimeout
	}

	return &tcpWatcher{
		ctx: context.Background(),

		config: config,
		item:   conf.Item,
	}, nil
}

func (w *tcpWatcher) Name() string {
	return w.item.Name()
}

func (w *tcpWatcher) Data() interface{} {
	return w.item.Data()
}

func (w *tcpWatcher) Helper() worker_queue.WatcherFunc {
	return func(i interface{}) (interface{}, error) {
		input := i.(*types.Endpoint)
		d := net.Dialer{Timeout: time.Duration(w.config.ConnectTimeout) * time.Millisecond}
		conn, err := d.Dial("tcp", input.ToString())
		if err != nil {
			return nil, fmt.Errorf("could not connect to server: %v", err)
		}

		defer conn.Close()

		return true, nil
	}
}

func (w *tcpWatcher) Notify(a interface{}, b interface{}) {
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
		return newTcpWatcher(conf)
	})
}

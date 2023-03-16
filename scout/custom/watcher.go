package custom

import (
	"fmt"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
)

const (
	modelName = "custom"
)

type customWatcher struct {
	config *CustomCheckerConfig
	item   scout.ScoutDelegate
}

func newCustomWatcher(conf *scout.WatcherConfig) (*customWatcher, error) {
	var config *CustomCheckerConfig = nil
	config, ok := conf.Data.(*CustomCheckerConfig)
	if !ok {
		return nil, fmt.Errorf("custom watcher invalid config: %v", conf)
	}

	if config.Probe == nil {
		return nil, fmt.Errorf("custom probe is nil config: %v", config)
	}

	return &customWatcher{
		config: config,
		item:   conf.Item,
	}, nil
}

func (w *customWatcher) Name() string {
	return w.item.Name()
}

func (w *customWatcher) Data() interface{} {
	return w.item.Data()
}

func (w *customWatcher) Helper() worker_queue.WatcherFunc {
	return w.config.Probe
}

func (w *customWatcher) Notify(a interface{}, b interface{}) {
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
		return newCustomWatcher(conf)
	})
}

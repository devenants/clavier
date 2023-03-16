package scout

import (
	"errors"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/types"
)

type ScoutDelegate interface {
	Name() string
	Data() interface{}
	Notify(*types.Endpoint, bool)
}

type WatcherConfig struct {
	Data interface{}
	Item ScoutDelegate
}

type ScoutWatcher interface {
	Name() string
	Data() interface{}
	Helper() worker_queue.WatcherFunc
	Notify(a interface{}, b interface{})
}

var (
	watcherFactoryByName = make(map[string]func(conf *WatcherConfig) (ScoutWatcher, error))
)

func Register(name string, factory func(conf *WatcherConfig) (ScoutWatcher, error)) {
	watcherFactoryByName[name] = factory
}

func ScoutWatcherCreate(name string, conf *WatcherConfig) (ScoutWatcher, error) {
	if f, ok := watcherFactoryByName[name]; ok {
		return f(conf)
	} else {
		return nil, errors.New("check model is not found")
	}
}

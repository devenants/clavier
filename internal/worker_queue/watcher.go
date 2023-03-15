package worker_queue

import (
	"fmt"
)

type WatcherFunc func(interface{}) (interface{}, error)

type WatcherDelegate interface {
	Name() string
	Data() interface{}
	Helper() WatcherFunc
	Notify(interface{}, interface{})
}

type Watcher struct {
	w WatcherDelegate
}

func NewWatcher(probe WatcherDelegate) (*Watcher, error) {
	if probe.Helper() == nil {
		return nil, fmt.Errorf("handler is nil tasklet: %v ", probe)
	}

	return &Watcher{
		w: probe,
	}, nil
}

func (w *Watcher) Name() string {
	return w.w.Name()
}

func (w *Watcher) Data() interface{} {
	return w.w.Data()
}

func (w *Watcher) Hander() HandlerFunc {
	return func(i interface{}) error {
		h := w.w.Helper()
		if h != nil {
			status, err := h(i)
			if err == nil {
				w.w.Notify(i, status)
			}
		}

		return nil
	}
}

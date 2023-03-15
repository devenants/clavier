package custom

import (
	"fmt"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
)

type checkWatcher struct {
	anchor *customChecker
	probe  scout.ScoutDelegate
}

func newCheckWatcher(anchor *customChecker, probe scout.ScoutDelegate) (*checkWatcher, error) {
	if anchor.config.Probe == nil {
		return nil, fmt.Errorf("handler is nil tasklet: %v config: %v", probe, anchor.config)
	}

	return &checkWatcher{
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
	watch := w.anchor.config.Probe
	return watch
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

package scout

import (
	"context"

	"github.com/devenants/clavier/internal/worker_queue"
)

type HelperConfig struct {
	Model        string
	Batches      int
	BatchTimeout int
	Data         interface{}
}

type CheckHelper struct {
	model string
	Data  interface{}

	wp *worker_queue.WorkerPool
}

func NewCheckHelper(conf *HelperConfig) (*CheckHelper, error) {
	batches := 0
	batchTimeout := 0
	if conf != nil {
		batches = conf.Batches
		batchTimeout = conf.BatchTimeout
	}

	poolConf := &worker_queue.PoolConfig{
		Batches:      batches,
		BatchTimeout: batchTimeout,
	}

	ctx := context.Background()
	p, err := worker_queue.NewWorkerPool(ctx, poolConf)
	if err != nil {
		return nil, err
	}

	return &CheckHelper{
		wp:    p,
		Data:  conf.Data,
		model: conf.Model,
	}, nil
}

func (t *CheckHelper) Model() string {
	return t.model
}

func (t *CheckHelper) Register(r ScoutDelegate) error {
	config := &WatcherConfig{
		Data: t.Data,
		Item: r,
	}
	ot, err := ScoutWatcherCreate(t.Model(), config)
	if err != nil {
		return err
	}

	w, err := worker_queue.NewWatcher(ot)
	if err != nil {
		return err
	}

	return t.wp.Register(w)
}

func (t *CheckHelper) UnRegister(r string) {
	t.wp.UnRegister(r)
}

func (t *CheckHelper) Done() {
	t.wp.Done()
}

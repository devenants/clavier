package custom

import (
	"context"
	"fmt"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
)

const (
	modelName = "custom"
)

type customChecker struct {
	wp     *worker_queue.WorkerPool
	config *CustomCheckerConfig
}

func NewCustomChecker(config *scout.ModelConfig) (*customChecker, error) {
	var conf *CustomCheckerConfig
	var ok bool
	if config.Data != nil {
		conf, ok = config.Data.(*CustomCheckerConfig)
		if !ok {
			return nil, fmt.Errorf("custom model data invalid %v", config)
		}
	}

	poolConf := &worker_queue.PoolConfig{
		Batches:      conf.Batches,
		BatchTimeout: conf.BatchTimeout,
	}

	ctx := context.Background()
	p, err := worker_queue.NewWorkerPool(ctx, poolConf)
	if err != nil {
		return nil, err
	}

	return &customChecker{
		wp:     p,
		config: conf,
	}, nil
}

func (t *customChecker) Model() string {
	return modelName
}

func (t *customChecker) Register(r scout.ScoutDelegate) error {
	ot, err := newCheckWatcher(t, r)
	if err != nil {
		return err
	}

	w, err := worker_queue.NewWatcher(ot)
	if err != nil {
		return err
	}

	return t.wp.Register(w)
}

func (t *customChecker) UnRegister(r string) {
	t.wp.UnRegister(r)
}

func (t *customChecker) Done() {
	t.wp.Done()
}

func init() {
	scout.Register(modelName, func(conf *scout.ModelConfig) (scout.CheckModel, error) {
		return NewCustomChecker(conf)
	})
}

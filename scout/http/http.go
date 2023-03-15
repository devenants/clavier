package http

import (
	"context"
	"fmt"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
)

const (
	modelName = "http"
)

type httpChecker struct {
	wp     *worker_queue.WorkerPool
	config *HttpCheckerConfig
}

func NewHttpChecker(config *scout.ModelConfig) (*httpChecker, error) {
	var conf *HttpCheckerConfig
	var ok bool
	if config.Data != nil {
		conf, ok = config.Data.(*HttpCheckerConfig)
		if !ok {
			return nil, fmt.Errorf("custom model data invalid %v", config)
		}
	}

	poolConf := &worker_queue.PoolConfig{
		Batches:      2,
		BatchTimeout: 5,

		CleanTimeout: 2,
		CleanMax:     2,
	}

	ctx := context.Background()
	p, err := worker_queue.NewWorkerPool(ctx, poolConf)
	if err != nil {
		return nil, err
	}

	return &httpChecker{
		wp:     p,
		config: conf,
	}, nil
}

func (t *httpChecker) Model() string {
	return modelName
}

func (t *httpChecker) Register(r scout.ScoutDelegate) error {
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

func (t *httpChecker) UnRegister(r string) {
	t.wp.UnRegister(r)
}

func (t *httpChecker) Done() {
	t.wp.Done()
}

func init() {
	scout.Register(modelName, func(conf *scout.ModelConfig) (scout.CheckModel, error) {
		return NewHttpChecker(conf)
	})
}

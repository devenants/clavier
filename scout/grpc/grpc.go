package grpc

import (
	"context"
	"fmt"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
)

const (
	modelName           = "grpc"
	defaultCheckTimeout = 5000
)

type grpcChecker struct {
	wp     *worker_queue.WorkerPool
	config *GrpcCheckerConfig
}

func NewGrpcChecker(config *scout.ModelConfig) (*grpcChecker, error) {
	var conf *GrpcCheckerConfig
	var ok bool
	if config.Data != nil {
		conf, ok = config.Data.(*GrpcCheckerConfig)
		if !ok {
			return nil, fmt.Errorf("custom model data invalid %v", config)
		}
		if conf.CheckTimeout == 0 {
			conf.CheckTimeout = defaultCheckTimeout
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

	return &grpcChecker{
		wp:     p,
		config: conf,
	}, nil
}

func (t *grpcChecker) Model() string {
	return modelName
}

func (t *grpcChecker) Register(r scout.ScoutDelegate) error {
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

func (t *grpcChecker) UnRegister(r string) {
	t.wp.UnRegister(r)
}

func (t *grpcChecker) Done() {
	t.wp.Done()
}

func init() {
	scout.Register(modelName, func(conf *scout.ModelConfig) (scout.CheckModel, error) {
		return NewGrpcChecker(conf)
	})
}

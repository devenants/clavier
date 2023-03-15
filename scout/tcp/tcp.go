package tcp

import (
	"context"
	"fmt"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
)

const (
	modelName             = "tcp"
	defaultConnectTimeout = 2000
)

type tcpChecker struct {
	wp     *worker_queue.WorkerPool
	config *TcpCheckerConfig
}

func NewTcpChecker(config *scout.ModelConfig) (*tcpChecker, error) {
	var conf *TcpCheckerConfig
	var ok bool
	if config.Data != nil {
		conf, ok = config.Data.(*TcpCheckerConfig)
		if !ok {
			return nil, fmt.Errorf("custom model data invalid %v", config)
		}
		if conf.ConnectTimeout == 0 {
			conf.ConnectTimeout = defaultConnectTimeout
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

	return &tcpChecker{
		wp:     p,
		config: conf,
	}, nil
}

func (t *tcpChecker) Model() string {
	return modelName
}

func (t *tcpChecker) Register(r scout.ScoutDelegate) error {
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

func (t *tcpChecker) UnRegister(r string) {
	t.wp.UnRegister(r)
}

func (t *tcpChecker) Done() {
	t.wp.Done()
}

func init() {
	scout.Register(modelName, func(conf *scout.ModelConfig) (scout.CheckModel, error) {
		return NewTcpChecker(conf)
	})
}

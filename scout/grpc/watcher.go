package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	modelName           = "grpc"
	defaultCheckTimeout = 5000
)

type grpcWatcher struct {
	config *GrpcCheckerConfig
	ctx    context.Context
	item   scout.ScoutDelegate
}

func newGrpcWatcher(conf *scout.WatcherConfig) (*grpcWatcher, error) {
	var config *GrpcCheckerConfig = nil
	config, ok := conf.Data.(*GrpcCheckerConfig)
	if !ok {
		return nil, fmt.Errorf("gRPC watcher invalid config: %v", conf)
	}

	if config.CheckTimeout == 0 {
		config.CheckTimeout = defaultCheckTimeout
	}

	return &grpcWatcher{
		config: config,
		ctx:    context.Background(),
		item:   conf.Item,
	}, nil
}

func (w *grpcWatcher) Name() string {
	return w.item.Name()
}

func (w *grpcWatcher) Data() interface{} {
	return w.item.Data()
}

func (w *grpcWatcher) Helper() worker_queue.WatcherFunc {
	return func(i interface{}) (interface{}, error) {
		input := i.(*types.Endpoint)
		conn, err := grpc.Dial(input.ToString(), w.config.DialOptions...)
		if err != nil {
			return nil, fmt.Errorf("gRPC health check failed on connect: %w", err)
		}
		defer conn.Close()

		healthClient := grpc_health_v1.NewHealthClient(conn)

		ctx, cancel := context.WithTimeout(w.ctx, time.Duration(w.config.CheckTimeout)*time.Millisecond)
		defer cancel()

		res, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
			Service: w.config.Service,
		})
		if err != nil {
			return nil, fmt.Errorf("gRPC health check failed on check call: %w", err)
		}

		if res.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
			return nil, fmt.Errorf("gRPC service reported as non-serving: %q", res.GetStatus().String())
		}

		return true, nil
	}
}

func (w *grpcWatcher) Notify(a interface{}, b interface{}) {
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
		return newGrpcWatcher(conf)
	})
}

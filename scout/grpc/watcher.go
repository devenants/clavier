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

type checkWatcher struct {
	ctx    context.Context
	anchor *grpcChecker
	probe  scout.ScoutDelegate
}

func newCheckWatcher(anchor *grpcChecker, probe scout.ScoutDelegate) (*checkWatcher, error) {
	return &checkWatcher{
		ctx:    context.Background(),
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
	return func(i interface{}) (interface{}, error) {
		input := i.(*types.Endpoint)
		conn, err := grpc.Dial(input.ToString(), w.anchor.config.DialOptions...)
		if err != nil {
			return nil, fmt.Errorf("gRPC health check failed on connect: %w", err)
		}
		defer conn.Close()

		healthClient := grpc_health_v1.NewHealthClient(conn)

		ctx, cancel := context.WithTimeout(w.ctx, time.Duration(w.anchor.config.CheckTimeout)*time.Millisecond)
		defer cancel()

		res, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
			Service: w.anchor.config.Service,
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

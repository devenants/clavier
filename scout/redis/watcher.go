package http

import (
	"context"
	"fmt"
	"time"

	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
	"github.com/redis/go-redis/v9"
)

const (
	modelName             = "redis"
	defaultRequestTimeout = 2000
)

type httpWatcher struct {
	config *RedisCheckerConfig
	ctx    context.Context
	item   scout.ScoutDelegate
}

func newHttpWatcher(conf *scout.WatcherConfig) (*httpWatcher, error) {
	var config *RedisCheckerConfig = nil
	config, ok := conf.Data.(*RedisCheckerConfig)
	if !ok {
		return nil, fmt.Errorf("http watcher invalid config: %v", conf)
	}

	if config.RequestTimeout == 0 {
		config.RequestTimeout = defaultRequestTimeout
	}

	return &httpWatcher{
		ctx:    context.Background(),
		config: config,
		item:   conf.Item,
	}, nil
}

func (w *httpWatcher) Name() string {
	return w.item.Name()
}

func (w *httpWatcher) Data() interface{} {
	return w.item.Data()
}

func (w *httpWatcher) Helper() worker_queue.WatcherFunc {
	return func(i interface{}) (interface{}, error) {
		input := i.(*types.Endpoint)
		dsn := fmt.Sprintf("redis://%s", input.ToString())
		redisOptions, _ := redis.ParseURL(dsn)

		ctx, cancel := context.WithTimeout(w.ctx, time.Duration(w.config.RequestTimeout)*time.Millisecond)
		defer cancel()

		rdb := redis.NewClient(redisOptions)
		defer rdb.Close()

		pong, err := rdb.Ping(ctx).Result()
		if err != nil {
			return false, fmt.Errorf("redis ping failed: %w", err)
		}

		if pong != "PONG" {
			return false, fmt.Errorf("unexpected response for redis ping: %q", pong)
		}

		return true, nil
	}
}

func (w *httpWatcher) Notify(a interface{}, b interface{}) {
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
		return newHttpWatcher(conf)
	})
}

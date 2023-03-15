package custom

import (
	"github.com/devenants/clavier/internal/worker_queue"
)

type CustomCheckerConfig struct {
	Probe worker_queue.WatcherFunc
}

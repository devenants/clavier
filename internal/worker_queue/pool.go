package worker_queue

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type PoolConfig struct {
	Batches      int
	BatchTimeout int

	CleanTimeout int32
	CleanMax     int32
}

type WorkerPool struct {
	ctx context.Context

	idx    int64
	config *PoolConfig

	workerMu sync.RWMutex
	queue    map[string]*Worker

	taskletWorkerMu    sync.RWMutex
	taskletWorkerQueue map[string]*Worker

	taskletMu    sync.RWMutex
	taskletQueue map[string]Tasklet
}

func NewWorkerPool(ctx context.Context, config *PoolConfig) (*WorkerPool, error) {
	p := &WorkerPool{
		ctx:                ctx,
		config:             config,
		queue:              make(map[string]*Worker),
		taskletWorkerQueue: make(map[string]*Worker),
		taskletQueue:       make(map[string]Tasklet),
	}

	go p.Run()

	return p, nil
}

func (p *WorkerPool) Done() {
	for _, v := range p.queue {
		v.Done()
	}

	p.ctx.Done()
}

func (p *WorkerPool) TaskletCount() int {
	p.workerMu.Lock()
	defer p.workerMu.Unlock()

	return len(p.taskletQueue)
}

func (p *WorkerPool) WorkerCount() int {
	p.workerMu.Lock()
	defer p.workerMu.Unlock()

	return len(p.queue)
}

func (p *WorkerPool) bind(r Tasklet, w *Worker) {
	p.taskletWorkerMu.Lock()
	p.taskletWorkerQueue[r.Name()] = w
	p.taskletWorkerMu.Unlock()

	p.taskletMu.Lock()
	p.taskletQueue[r.Name()] = r
	p.taskletMu.Unlock()
}

func (p *WorkerPool) UnRegister(r string) {
	p.taskletWorkerMu.Lock()
	if v, ok := p.taskletWorkerQueue[r]; ok {
		v.DelTasklet(r)
		delete(p.taskletWorkerQueue, r)
	}
	p.taskletWorkerMu.Unlock()

	p.taskletMu.Lock()
	delete(p.taskletQueue, r)
	p.taskletMu.Unlock()
}

func (p *WorkerPool) Register(r Tasklet) error {
	for _, v := range p.queue {
		if v.RequestAvailable() {
			ret := v.AddTasklet(r)
			if ret == WorkerNormal || ret == WorkerReplace {
				p.bind(r, v)
				return nil
			}
		}
	}

	config := &WorkerConfig{
		batches:      p.config.Batches,
		batchTimeout: p.config.BatchTimeout,
	}

	name := strconv.FormatInt(p.idx, 10)
	atomic.AddInt64(&p.idx, 1)

	w, err := NewWorker(name, config)
	if err != nil {
		return err
	}

	ret := w.AddTasklet(r)
	if ret != WorkerNormal && ret != WorkerReplace {
		return fmt.Errorf("requeust add failed ret: %v w: %v r: %v", ret, w, r)
	}

	p.bind(r, w)

	p.workerMu.Lock()
	p.queue[name] = w
	p.workerMu.Unlock()

	return nil
}

func (p *WorkerPool) Run() error {
	ticker := time.NewTicker(time.Second * time.Duration(p.config.CleanTimeout))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, v := range p.queue {
				reqCount := v.TaskletCount()
				if reqCount == 0 {
					count := v.DiscardCount()
					if (count + 1) >= p.config.CleanMax {
						v.Done()
						p.workerMu.Lock()
						delete(p.queue, v.Name())
						p.workerMu.Unlock()
					} else {
						v.DiscardInc()
					}
				}

			}
		case <-p.ctx.Done():
			return p.ctx.Err()
		}
	}
}

package worker_queue

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	WorkerFull = iota
	WorkerNotFound
	WorkerNormal
	WorkerReplace
)

type HandlerFunc func(interface{}) error

type Tasklet interface {
	Name() string
	Data() interface{}
	Hander() HandlerFunc
}

type WorkerConfig struct {
	batches      int
	batchTimeout int
}

type Worker struct {
	status      int32
	name        string
	ctx         context.Context
	waitMu      sync.RWMutex
	watingQueue map[string]Tasklet

	batches      int
	batchTimeout int

	discardCount int32
}

func NewWorker(name string, config *WorkerConfig) (*Worker, error) {
	w := &Worker{
		name:         name,
		ctx:          context.Background(),
		watingQueue:  make(map[string]Tasklet),
		batches:      config.batches,
		batchTimeout: config.batchTimeout,
	}

	go w.Run()

	return w, nil
}

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) Done() {
	atomic.AddInt32(&w.status, 1)
	w.ctx.Done()
}

func (w *Worker) Status() bool {
	status := atomic.LoadInt32(&w.status)
	return status == 0
}

func (w *Worker) DiscardInc() {
	atomic.AddInt32(&w.discardCount, 1)
}

func (w *Worker) DiscardCount() int32 {
	return atomic.LoadInt32(&w.discardCount)
}

func (w *Worker) RequestAvailable() bool {
	w.waitMu.Lock()
	defer w.waitMu.Unlock()

	return len(w.watingQueue) < w.batches
}

func (w *Worker) TaskletCount() int {
	w.waitMu.Lock()
	defer w.waitMu.Unlock()

	return len(w.watingQueue)
}

func (w *Worker) AddTasklet(r Tasklet) int {
	w.waitMu.Lock()
	defer w.waitMu.Unlock()

	if len(w.watingQueue) >= w.batches {
		return WorkerFull
	}

	if _, ok := w.watingQueue[r.Name()]; ok {
		w.watingQueue[r.Name()] = r
		return WorkerReplace
	}

	w.watingQueue[r.Name()] = r
	return WorkerNormal
}

func (w *Worker) DelTasklet(name string) int {
	w.waitMu.Lock()
	defer w.waitMu.Unlock()

	if _, ok := w.watingQueue[name]; !ok {
		return WorkerNotFound
	}

	delete(w.watingQueue, name)

	return WorkerNormal
}

func (w *Worker) Run() error {
	ticker := time.NewTicker(time.Second * time.Duration(w.batchTimeout))
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return w.ctx.Err()
		case <-ticker.C:
			w.waitMu.Lock()
			pq := make([]Tasklet, 0)
			for _, v := range w.watingQueue {
				pq = append(pq, v)
			}
			w.waitMu.Unlock()

			for _, v := range pq {
				h := v.Hander()
				if h != nil {
					err := h(v.Data())
					if err != nil {
						fmt.Printf("handler failed %v %v", err, v.Data())
					}
				}
			}
		}
	}
}

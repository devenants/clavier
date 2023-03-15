package clavier

import (
	"context"
	"fmt"
	"sync"

	_ "github.com/devenants/clavier/discovery/dns"
	_ "github.com/devenants/clavier/filter/default"
	_ "github.com/devenants/clavier/filter/round_robin"
	"github.com/devenants/clavier/types"
)

type Clavier struct {
	ctx context.Context

	mu        sync.RWMutex
	listeners map[string]*Listener
}

func NewClavier(ctx context.Context) (*Clavier, error) {
	c := &Clavier{
		listeners: make(map[string]*Listener),
		ctx:       ctx,
	}

	go c.Run()

	return c, nil
}

func (r *Clavier) ListenerCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.listeners)
}

func (r *Clavier) AddListener(name string, config *ListenerConfig) (*Listener, error) {
	ctx := context.Background()
	lis, err := NewListener(ctx, config, nil)
	if err != nil {
		return nil, err
	}

	go lis.Run()

	r.mu.Lock()
	defer r.mu.Unlock()
	r.listeners[name] = lis

	return lis, nil
}

func (r *Clavier) DelListener(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	lis, ok := r.listeners[name]
	if ok {
		lis.ctx.Done()
	}

	delete(r.listeners, name)
}

func (r *Clavier) ListEndpoints(name string) ([]*types.Endpoint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	lis, ok := r.listeners[name]
	if ok {
		return lis.ListEndpoints(), nil
	}

	return nil, fmt.Errorf("listener not found: %v", name)
}

func (r *Clavier) GetEndpoint(name string) (*types.Endpoint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	lis, ok := r.listeners[name]
	if ok {
		return lis.GetEndpoint()
	}

	return nil, fmt.Errorf("listener not found")
}

func (r *Clavier) Run() error {
	<-r.ctx.Done()
	r.mu.Lock()
	for _, v := range r.listeners {
		v.ctx.Done()
	}
	r.mu.Unlock()

	return nil
}

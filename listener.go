package clavier

import (
	"context"
	"fmt"
	"time"

	_ "github.com/devenants/clavier/discovery/dns"
	_ "github.com/devenants/clavier/filter/default"
	_ "github.com/devenants/clavier/filter/round_robin"
	_ "github.com/devenants/clavier/scout/custom"
	"github.com/devenants/clavier/types"

	"github.com/devenants/clavier/discovery"
)

type DiscoveryConfig struct {
	Model  string
	Config discovery.ModelConfig
}

type ListenerConfig struct {
	Group string
	Items []*types.Endpoint

	Discovery *DiscoveryConfig

	Entry *EntryManagerConfig
}

type Listener struct {
	ctx           context.Context
	m             *ListenerConfig
	privateLookup discovery.LookupFunc
	dis           discovery.DiscoveryModel
	dm            *entryManger
}

func NewListener(ctx context.Context, config *ListenerConfig, lookup func(interface{}) ([]*types.Endpoint, error)) (*Listener, error) {
	dis, err := discovery.DisModelCreate(config.Discovery.Model, &config.Discovery.Config)
	if err != nil {
		return nil, err
	}

	dm, err := NewDstManger(config.Entry)
	if err != nil {
		return nil, err
	}

	l := &Listener{
		ctx:           ctx,
		m:             config,
		privateLookup: lookup,
		dis:           dis,
		dm:            dm,
	}

	go l.Run()

	return l, nil
}

func (l *Listener) ListEndpoints() []*types.Endpoint {
	return l.dm.dsts()
}

func (l *Listener) GetEndpoint() (*types.Endpoint, error) {
	return l.dm.dst()
}

func (l *Listener) lookup() error {
	var dst []*types.Endpoint
	var err error
	if len(l.m.Group) > 0 {
		if l.privateLookup != nil {
			dst, err = l.privateLookup(l.m.Group)
			if err != nil {
				return fmt.Errorf("lookup failed %v %v", l.m.Group, err)
			}
		} else {
			dst, err = l.dis.Lookup(l.m.Group, nil)
			if err != nil {
				return fmt.Errorf("lookup failed %v %v", l.m.Group, err)
			}
		}
	} else {
		dst = l.m.Items
	}

	l.dm.enqueue(dst)

	return nil
}

func (l *Listener) Run() error {
	for {
		select {
		case <-l.ctx.Done():
			return l.ctx.Err()
		default:
			err := l.lookup()
			if err != nil {
				return fmt.Errorf("lookup failed %v %v", l.m.Group, err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

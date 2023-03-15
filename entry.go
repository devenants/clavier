package clavier

import (
	"fmt"
	"sync"

	_ "github.com/devenants/clavier/discovery/dns"
	"github.com/devenants/clavier/filter"
	_ "github.com/devenants/clavier/filter/default"
	_ "github.com/devenants/clavier/filter/round_robin"
	"github.com/devenants/clavier/internal/worker_queue"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
)

type FilterMangerConfig struct {
	Model      string
	StatusJump bool
	Config     filter.ModelConfig
}

type ScoutMangerConfig struct {
	Model  string
	Config scout.ModelConfig
}

type EntryManagerConfig struct {
	FilterConfig *FilterMangerConfig
	ScoutConfig  *ScoutMangerConfig
}

type entryManger struct {
	cacheMu   sync.RWMutex
	cachedDst []*types.Endpoint

	f filter.FilterModel

	probe   worker_queue.WatcherFunc
	checker scout.CheckModel

	entryMu       sync.RWMutex
	launchedQueue map[string]*probeEntry

	config *EntryManagerConfig
}

func NewDstManger(emc *EntryManagerConfig) (*entryManger, error) {
	if len(emc.FilterConfig.Model) == 0 {
		emc.FilterConfig.Model = "default"
	}
	f, err := filter.FilterModelCreate(emc.FilterConfig.Model, &emc.FilterConfig.Config)
	if err != nil {
		return nil, err
	}

	var s scout.CheckModel
	if emc.ScoutConfig.Model != "none" {
		s, err = scout.CheckModelCreate(emc.ScoutConfig.Model, &emc.ScoutConfig.Config)
		if err != nil {
			return nil, err
		}
	}

	return &entryManger{
		f:             f,
		cachedDst:     make([]*types.Endpoint, 0),
		launchedQueue: make(map[string]*probeEntry),
		checker:       s,
		config:        emc,
	}, nil
}

func (m *entryManger) update(host *types.Endpoint, status bool) {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	m.launchedQueue[host.ToString()].entry.Status = status
}

func (m *entryManger) launch(entries []*types.Endpoint) {
	m.entryMu.Lock()
	defer m.entryMu.Unlock()

	if m.checker == nil {
		return
	}

	for _, e := range entries {
		if _, ok := m.launchedQueue[e.ToString()]; !ok {
			ne := &probeEntry{
				entry: e,
				m:     m,
				probe: m.probe,
			}

			m.checker.Register(ne)
			m.launchedQueue[e.ToString()] = ne
		}
	}

	//TODO clean
}

func (m *entryManger) enqueue(dsts []*types.Endpoint) {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	m.cachedDst = dsts
	m.launch(dsts)
}

func (m *entryManger) dsts() []*types.Endpoint {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	return m.cachedDst

}

func (m *entryManger) dst() (*types.Endpoint, error) {
	if m.checker == nil {
		m.cacheMu.Lock()
		defer m.cacheMu.Unlock()

		h, err := m.f.Shuffle(m.cachedDst, m.config.FilterConfig.StatusJump)
		if err != nil {
			return nil, err
		}

		return h, nil
	}

	m.entryMu.Lock()
	defer m.entryMu.Unlock()

	waitingList := make([]*types.Endpoint, 0)
	for _, v := range m.launchedQueue {
		waitingList = append(waitingList, v.entry)
	}

	h, err := m.f.Shuffle(waitingList, m.config.FilterConfig.StatusJump)
	if err != nil {
		return nil, err
	}
	return h, nil
}

type probeEntry struct {
	entry *types.Endpoint
	probe worker_queue.WatcherFunc
	m     *entryManger
}

func (t *probeEntry) Name() string {
	return fmt.Sprintf("%s-%s", t.entry.Host, t.entry.Port)
}

func (t *probeEntry) Data() interface{} {
	return t.entry
}

func (t *probeEntry) Notify(host *types.Endpoint, status bool) {
	t.m.update(host, status)
}

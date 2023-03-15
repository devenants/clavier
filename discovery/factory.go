package discovery

import (
	"errors"

	"github.com/devenants/clavier/types"
)

type DiscoveryModel interface {
	Model() string
	Lookup(name string, option interface{}) ([]*types.Endpoint, error)
}

var (
	modelFactoryByName = make(map[string]func(conf *ModelConfig) (DiscoveryModel, error))
)

func Register(name string, factory func(conf *ModelConfig) (DiscoveryModel, error)) {
	modelFactoryByName[name] = factory
}

func DisModelCreate(name string, conf *ModelConfig) (DiscoveryModel, error) {
	if f, ok := modelFactoryByName[name]; ok {
		return f(conf)
	} else {
		return nil, errors.New("discovery model is not found")
	}
}

type ModelConfig struct {
	Data interface{}
}

type LookupFunc func(interface{}) ([]*types.Endpoint, error)

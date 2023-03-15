package filter

import (
	"errors"

	"github.com/devenants/clavier/types"
)

type FilterModel interface {
	Model() string
	Shuffle([]*types.Endpoint, interface{}) (*types.Endpoint, error)
}

var (
	modelFactoryByName = make(map[string]func(conf *ModelConfig) (FilterModel, error))
)

func Register(name string, factory func(conf *ModelConfig) (FilterModel, error)) {
	modelFactoryByName[name] = factory
}

func FilterModelCreate(name string, conf *ModelConfig) (FilterModel, error) {
	if f, ok := modelFactoryByName[name]; ok {
		return f(conf)
	} else {
		return nil, errors.New("filter model is not found")
	}
}

type ModelConfig struct {
	Data interface{}
}

type ShuffleFunc func(interface{}) ([]types.Endpoint, error)

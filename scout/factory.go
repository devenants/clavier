package scout

import (
	"errors"

	"github.com/devenants/clavier/types"
)

type CheckModel interface {
	Model() string

	Register(ScoutDelegate) error
	UnRegister(string)

	Done()
}

var (
	modelFactoryByName = make(map[string]func(conf *ModelConfig) (CheckModel, error))
)

func Register(name string, factory func(conf *ModelConfig) (CheckModel, error)) {
	modelFactoryByName[name] = factory
}

func CheckModelCreate(name string, conf *ModelConfig) (CheckModel, error) {
	if f, ok := modelFactoryByName[name]; ok {
		return f(conf)
	} else {
		return nil, errors.New("check model is not found")
	}
}

type ModelConfig struct {
	Data interface{}
}

type ScoutDelegate interface {
	Name() string
	Data() interface{}
	Notify(*types.Endpoint, bool)
}

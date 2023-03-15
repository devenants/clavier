package headfirst

import (
	"fmt"

	"github.com/devenants/clavier/filter"
	"github.com/devenants/clavier/types"
)

const (
	modelName = "default"
)

type DefaultFilter struct {
}

func NewDefaultFilter(_ *filter.ModelConfig) (*DefaultFilter, error) {
	return &DefaultFilter{}, nil
}

func (r *DefaultFilter) Model() string {
	return modelName
}

func (r *DefaultFilter) Shuffle(items []*types.Endpoint, option interface{}) (*types.Endpoint, error) {
	if len(items) > 0 {
		return items[0], nil
	}
	return nil, fmt.Errorf("items is nil")
}

func init() {
	filter.Register(modelName, func(conf *filter.ModelConfig) (filter.FilterModel, error) {
		return NewDefaultFilter(conf)
	})
}

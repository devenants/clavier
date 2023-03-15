package round_robin

import (
	"fmt"
	"sort"
	"strings"

	"github.com/devenants/clavier/filter"
	"github.com/devenants/clavier/types"
)

const (
	modelName = "round-robin"
)

type RrFilter struct {
	idx int
}

func NewRrFilter(_ *filter.ModelConfig) (*RrFilter, error) {
	return &RrFilter{
		idx: 0,
	}, nil
}

func (r *RrFilter) Model() string {
	return modelName
}

func (r *RrFilter) Shuffle(items []*types.Endpoint, option interface{}) (*types.Endpoint, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("items is nil")
	}

	sort.SliceStable(items, func(i, j int) bool {
		return strings.Compare(items[i].Host, items[j].Host) < 0
	})

	jump := option.(bool)

	count := 0
	for {
		if count >= len(items) {
			break
		}

		if r.idx >= len(items) {
			r.idx = 0
		}

		if items[r.idx].Status || jump {
			resp := items[r.idx]
			r.idx += 1
			return resp, nil
		}

		count += 1
		r.idx += 1
	}

	return nil, fmt.Errorf("health endpoint is zero")
}

func init() {
	filter.Register(modelName, func(conf *filter.ModelConfig) (filter.FilterModel, error) {
		return NewRrFilter(conf)
	})
}

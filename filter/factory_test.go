package filter

import (
	"testing"

	"github.com/devenants/clavier/types"
	"github.com/stretchr/testify/require"
)

type testFilterFactory struct{}

func (f *testFilterFactory) Model() string {
	return ""
}

func (f *testFilterFactory) Shuffle([]*types.Endpoint, interface{}) (*types.Endpoint, error) {
	return &types.Endpoint{}, nil
}

func testFilterFactoryCreator(conf *ModelConfig) (FilterModel, error) {
	return &testFilterFactory{}, nil
}

func TestCreateFilterFactory(t *testing.T) {
	name := "test"
	Register(name, testFilterFactoryCreator)

	if m, err := FilterModelCreate(name, &ModelConfig{
		Data: nil,
	}); err != nil {
		t.Error(err)
		require.Equal(t, m.Model(), name, "")
	}
}

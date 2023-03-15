package discovery

import (
	"testing"

	"github.com/devenants/clavier/types"
	"github.com/stretchr/testify/require"
)

type testDiscoveryFactory struct{}

func (f *testDiscoveryFactory) Model() string {
	return ""
}

func (f *testDiscoveryFactory) Lookup(name string, option interface{}) ([]*types.Endpoint, error) {
	ends := make([]*types.Endpoint, 0)
	return ends, nil
}

func testDiscoveryFactoryCreator(conf *ModelConfig) (DiscoveryModel, error) {
	return &testDiscoveryFactory{}, nil
}

func TestCreateDiscoveryFactory(t *testing.T) {
	name := "test"
	Register(name, testDiscoveryFactoryCreator)

	if m, err := DisModelCreate(name, &ModelConfig{
		Data: nil,
	}); err != nil {
		t.Error(err)
		require.Equal(t, m.Model(), name, "")
	}
}

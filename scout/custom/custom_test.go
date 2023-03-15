package custom

import (
	"testing"

	"github.com/devenants/clavier/scout"
	"github.com/stretchr/testify/require"
)

func TestCustomWatcherTest(t *testing.T) {
	sc := &scout.ModelConfig{
		Data: &CustomCheckerConfig{
			Probe: func(_ interface{}) (interface{}, error) {
				return true, nil
			},
		},
	}

	c, err := NewCustomChecker(sc)
	require.Equal(t, err, nil, "")

	n := c.Model()
	require.Equal(t, n, modelName, "")

	e := &testEntry{
		name: "192.168.11.2",
		idx:  0,
	}
	err = c.Register(e)
	require.Equal(t, err, nil, "")

	c.UnRegister(e.name)

	c.Done()
}

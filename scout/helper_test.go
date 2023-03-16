package scout

import (
	"testing"

	"github.com/devenants/clavier/types"
	"github.com/stretchr/testify/require"
)

func TestCheckHelper(t *testing.T) {
	config := &HelperConfig{
		Model: "test",
	}

	ch, err := NewCheckHelper(config)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, ch, nil, "")

	require.Equal(t, ch.Model(), "test", "")

	e := &testEntry{
		addr: types.Endpoint{
			Host: "192.168.11.2",
			Port: "80",
		},
		name: "192.168.11.2",
		idx:  0,
	}

	err = ch.Register(e)
	require.Equal(t, err, nil, "")

	ch.UnRegister(e.name)

	ch.Done()
}

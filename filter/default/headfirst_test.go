package headfirst

import (
	"testing"

	"github.com/devenants/clavier/types"
	"github.com/stretchr/testify/require"
)

func TestWatcherCreateTest(t *testing.T) {
	f, err := NewDefaultFilter(nil)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, f, nil, "")

	require.Equal(t, f.Model(), modelName, "")

	a := []*types.Endpoint{
		{
			Host: "10.10.10.10",
			Port: "8080",
		},
		{
			Host: "10.10.10.11",
			Port: "8080",
		},
	}

	d, err := f.Shuffle(a, nil)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, d, nil, "")

}

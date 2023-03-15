package round_robin

import (
	"testing"

	"github.com/devenants/clavier/types"
	"github.com/stretchr/testify/require"
)

var (
	rf  *RrFilter
	err error
)

func TestRrFilterCreateTest(t *testing.T) {
	rf, err = NewRrFilter(nil)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, rf, nil, "")

	t.Run("testRrTest", testRrTest)
}

func testRrTest(t *testing.T) {
	require.Equal(t, rf.Model(), modelName, "")

	a := []*types.Endpoint{
		{
			Host:   "10.10.10.10",
			Port:   "8080",
			Status: true,
		},
		{
			Host:   "10.10.10.11",
			Port:   "8080",
			Status: true,
		},
	}

	d1, err := rf.Shuffle(a, false)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, d1, nil, "")

	d2, err := rf.Shuffle(a, false)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, d2, nil, "")

	require.NotEqual(t, d1, d2, "")

	b := []*types.Endpoint{
		{
			Host:   "10.10.10.10",
			Port:   "8080",
			Status: false,
		},
		{
			Host:   "10.10.10.11",
			Port:   "8080",
			Status: true,
		},
	}

	e1, err := rf.Shuffle(b, false)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, e1, nil, "")

	e2, err := rf.Shuffle(b, false)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, e2, nil, "")

	require.Equal(t, e1, e2, "")

}

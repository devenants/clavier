package custom

import (
	"testing"

	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
	"github.com/stretchr/testify/require"
)

type testEntry struct {
	name string
	idx  int
}

func (t *testEntry) Name() string {
	return t.name
}

func (t *testEntry) Data() interface{} {
	return t.name
}

func (t *testEntry) Notify(host *types.Endpoint, status bool) {
	t.idx += 1
}

var (
	cw  *customWatcher
	err error
)

func TestCheckWatcherTest(t *testing.T) {
	e := &testEntry{
		name: "192.168.11.2",
		idx:  0,
	}

	sc := &scout.WatcherConfig{
		Item: e,
		Data: &CustomCheckerConfig{
			Probe: func(_ interface{}) (interface{}, error) {
				return true, nil
			},
		},
	}

	cw, err = newCustomWatcher(sc)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, cw, nil, "")

	n := cw.Name()
	require.Equal(t, n, e.name, "")

	d := cw.Data()
	require.Equal(t, d, e.name, "")

	h := cw.Helper()
	ret, err := h(e)
	require.Equal(t, err, nil, "")
	require.Equal(t, ret.(bool), true, "")

	cw.Notify(&types.Endpoint{}, false)
	require.Equal(t, e.idx, 1, "")
}

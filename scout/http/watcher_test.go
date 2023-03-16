package http

import (
	"testing"

	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/types"
	"github.com/stretchr/testify/require"
)

type testEntry struct {
	addr types.Endpoint
	name string
	idx  int
}

func (t *testEntry) Name() string {
	return t.name
}

func (t *testEntry) Data() interface{} {
	return &t.addr
}

func (t *testEntry) Notify(host *types.Endpoint, status bool) {
	t.idx += 1
}

var (
	cw *httpWatcher
)

func TestHttpCheckWatcher(t *testing.T) {
	e := &testEntry{
		addr: types.Endpoint{
			Host: "192.168.11.2",
			Port: "80",
		},
		name: "192.168.11.2",
		idx:  0,
	}

	sc := &scout.WatcherConfig{
		Item: e,
		Data: &HttpCheckerConfig{
			RequestTimeout: 1000,
			URL:            "/",
			Method:         "GET",
		},
	}

	cw, err := newHttpWatcher(sc)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, cw, nil, "")

	n := cw.Name()
	require.Equal(t, n, e.name, "")

	d := cw.Data()
	require.Equal(t, d.(*types.Endpoint).Host, e.addr.Host, "")

	h := cw.Helper()
	ret, err := h(d)
	require.Equal(t, err, nil, "")
	require.Equal(t, ret.(bool), true, "")

	cw.Notify(&types.Endpoint{}, false)
	require.Equal(t, e.idx, 1, "")
}

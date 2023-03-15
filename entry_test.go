package clavier

import (
	"fmt"
	"testing"

	_ "github.com/devenants/clavier/discovery/dns"
	"github.com/devenants/clavier/filter"
	_ "github.com/devenants/clavier/filter/default"
	_ "github.com/devenants/clavier/filter/round_robin"
	"github.com/devenants/clavier/scout"
	"github.com/devenants/clavier/scout/custom"
	"github.com/devenants/clavier/types"
	"github.com/stretchr/testify/require"
)

var (
	dm  *entryManger
	err error
)

func TestEntrymanager(t *testing.T) {
	config := &EntryManagerConfig{
		FilterConfig: &FilterMangerConfig{
			Model:      "round-robin",
			StatusJump: true,
			Config: filter.ModelConfig{
				Data: nil,
			},
		},
		ScoutConfig: &ScoutMangerConfig{
			Model: "custom",
			Config: scout.ModelConfig{
				Data: &custom.CustomCheckerConfig{
					Probe: func(_ interface{}) (interface{}, error) {
						return true, nil
					},
				},
			},
		},
	}
	dm, err = NewDstManger(config)
	if err != nil {
		t.Error(err)
	}

	t.Run("testEntryLauch", testEntryLauch)
	t.Run("testEntryDst", testEntryDst)
	t.Run("testProbeEntry", testProbeEntry)
}

func testEntryLauch(t *testing.T) {
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

	dm.launch(a)

	require.Equal(t, len(dm.launchedQueue), 2, "")
}

func testEntryDst(t *testing.T) {
	a := []*types.Endpoint{
		{
			Host:   "192.10.10.10",
			Port:   "8080",
			Status: true,
		},
		{
			Host:   "202.16.10.11",
			Port:   "8080",
			Status: true,
		},
	}

	dm.enqueue(a)
	require.Equal(t, len(dm.cachedDst), 2, "")

	b := types.Endpoint{
		Host: "192.10.10.10",
		Port: "8080",
	}
	dm.update(&b, false)

	c := types.Endpoint{
		Host: "202.16.10.11",
		Port: "8080",
	}
	dm.update(&c, false)

	for _, item := range dm.cachedDst {
		require.Equal(t, item.Status, false, "")
	}

	d, err := dm.dst()
	require.Equal(t, err, nil, "")
	require.NotEqual(t, d, nil, "")

	ds := dm.dsts()
	require.Equal(t, len(ds), 2, "")
}

func testProbeEntry(t *testing.T) {
	a := types.Endpoint{
		Host: "100.64.10.2",
		Port: "443",
	}

	ne := &probeEntry{
		entry: &a,
		m:     dm,
		probe: nil,
	}

	require.Equal(t, ne.Name(), fmt.Sprintf("%s-%s", a.Host, a.Port), "")
	require.Equal(t, ne.Data(), &a, "")

	b := types.Endpoint{
		Host: "192.10.10.10",
		Port: "8080",
	}

	ne.Notify(&b, true)

	c := types.Endpoint{
		Host: "202.16.10.11",
		Port: "8080",
	}
	ne.Notify(&c, true)

	for _, item := range dm.launchedQueue {
		require.Equal(t, item.entry.Status, true, "")
	}
}

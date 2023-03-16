package scout

import (
	"testing"

	"github.com/devenants/clavier/internal/worker_queue"
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

type testScoutFactory struct{}

func (f *testScoutFactory) Name() string {
	return ""
}

func (f *testScoutFactory) Data() interface{} {
	return nil
}

func (f *testScoutFactory) Helper() worker_queue.WatcherFunc {
	return func(i interface{}) (interface{}, error) {
		return true, nil
	}
}

func (f *testScoutFactory) Notify(a interface{}, b interface{}) {
}

func testFactoryCreator(conf *WatcherConfig) (ScoutWatcher, error) {
	return &testScoutFactory{}, nil
}

func TestCreateTestFactory(t *testing.T) {
	e := &testEntry{
		addr: types.Endpoint{
			Host: "192.168.11.2",
			Port: "80",
		},
		name: "192.168.11.2",
		idx:  0,
	}

	name := "test"
	Register(name, testFactoryCreator)

	if m, err := ScoutWatcherCreate(name, &WatcherConfig{
		Data: nil,
		Item: e,
	}); err != nil {
		t.Error(err)
		require.Equal(t, m.Name(), e.name, "")
	}
}

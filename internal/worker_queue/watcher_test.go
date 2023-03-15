package worker_queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Item struct {
	name string
	idx  int
}

func (t *Item) Name() string {
	return t.name
}

func (t *Item) Data() interface{} {
	return t.name
}

func (t *Item) Helper() WatcherFunc {
	return func(i interface{}) (interface{}, error) {
		t.idx += 1
		return true, nil
	}
}

func (t *Item) Notify(interface{}, interface{}) {
}

var (
	e   *Watcher
	err error
)

func TestWatcherCreateTest(t *testing.T) {
	req1 := &Item{
		name: "test1",
		idx:  0,
	}

	e, err = NewWatcher(req1)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, w, nil, "")

	name := e.Name()
	require.Equal(t, name, req1.name, "")

	data := e.Data()
	require.Equal(t, data.(string), req1.name, "")

	h := e.Hander()
	err := h(nil)
	require.Equal(t, err, nil, "")

	require.Equal(t, req1.idx, 1, "")
}

package worker_queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Task struct {
	name string
	idx  int
}

func (t *Task) Name() string {
	return t.name
}

func (t *Task) Data() interface{} {
	return t
}

func (t *Task) Hander() HandlerFunc {
	return func(i interface{}) error {
		t := i.(*Task)
		t.idx += 1
		return nil
	}
}

func (t *Task) Notify(interface{}, interface{}) {

}

var (
	w *Worker
)

func TestWorkerCreateTest(t *testing.T) {
	var err error
	config := &WorkerConfig{
		batches:      2,
		batchTimeout: 5,
	}

	w, err = NewWorker("test", config)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, w, nil, "")

	req1 := &Task{
		name: "test1",
		idx:  0,
	}

	ret := w.AddTasklet(req1)
	require.Equal(t, ret, WorkerNormal, "")

	t.Run("testRequestTest", testRequestTest)
	t.Run("testAvailableTest", testAvailableTest)
	t.Run("testStatusTest", testStatusTest)
}

func testRequestTest(t *testing.T) {
	req2 := &Task{
		name: "test2",
		idx:  0,
	}

	ret := w.AddTasklet(req2)
	require.Equal(t, ret, WorkerNormal, "")
	require.Equal(t, req2.idx, 0, "")

	time.Sleep(time.Duration(w.batchTimeout*2) * time.Second)
	require.NotEqual(t, req2.idx, 0, "")

	req3 := &Task{
		name: "test3",
		idx:  0,
	}

	ret = w.AddTasklet(req3)
	require.Equal(t, ret, WorkerFull, "")

	ava := w.RequestAvailable()
	require.Equal(t, ava, false, "")

	ret = w.DelTasklet(req2.name)
	require.Equal(t, ret, WorkerNormal, "")

	count := w.TaskletCount()
	require.Equal(t, count, 1, "")

	ava = w.RequestAvailable()
	require.Equal(t, ava, true, "")

}

func testAvailableTest(t *testing.T) {
	count := w.DiscardCount()
	require.Equal(t, count, int32(0), "")

	w.DiscardInc()

	count = w.DiscardCount()
	require.Equal(t, count, int32(1), "")
}

func testStatusTest(t *testing.T) {
	status := w.Status()
	require.Equal(t, status, true, "")

	w.Done()

	status = w.Status()
	require.Equal(t, status, false, "")
}

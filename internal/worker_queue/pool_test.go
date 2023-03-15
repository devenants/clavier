package worker_queue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	p *WorkerPool
)

func TestPoolCreateTest(t *testing.T) {
	var err error
	config := &PoolConfig{
		Batches:      2,
		BatchTimeout: 5,

		CleanTimeout: 2,
		CleanMax:     2,
	}

	ctx := context.Background()
	p, err = NewWorkerPool(ctx, config)
	require.Equal(t, err, nil, "")
	require.NotEqual(t, p, nil, "")

	t.Run("testRegisterTest", testRegisterTest)
}

func testRegisterTest(t *testing.T) {
	req1 := &Task{
		name: "test1",
		idx:  0,
	}
	err := p.Register(req1)
	require.Equal(t, err, nil, "")

	time.Sleep(time.Duration(p.config.BatchTimeout*2) * time.Second)
	require.NotEqual(t, req1.idx, 0, "")

	count := p.WorkerCount()
	require.Equal(t, count, 1, "")

	req2 := &Task{
		name: "test2",
		idx:  0,
	}

	err = p.Register(req2)
	require.Equal(t, err, nil, "")

	count = p.WorkerCount()
	require.Equal(t, count, 1, "")

	count = p.TaskletCount()
	require.Equal(t, count, 2, "")

	req3 := &Task{
		name: "test3",
		idx:  0,
	}

	err = p.Register(req3)
	require.Equal(t, err, nil, "")

	count = p.WorkerCount()
	require.Equal(t, count, 2, "")

	count = p.TaskletCount()
	require.Equal(t, count, 3, "")

	p.UnRegister(req3.name)
	count = p.TaskletCount()
	require.Equal(t, count, 2, "")

	p.UnRegister(req2.name)
	count = p.TaskletCount()
	require.Equal(t, count, 1, "")

	delay := p.config.CleanTimeout * p.config.CleanMax
	time.Sleep(time.Duration(delay+1) * time.Second)

	count = p.WorkerCount()
	require.Equal(t, count, 1, "")

	p.UnRegister(req1.name)

	time.Sleep(time.Duration(delay+1) * time.Second)

	count = p.TaskletCount()
	require.Equal(t, count, 0, "")

	count = p.WorkerCount()
	require.Equal(t, count, 0, "")

	req4 := &Task{
		name: "test4",
		idx:  0,
	}

	err = p.Register(req4)
	require.Equal(t, err, nil, "")

	count = p.TaskletCount()
	require.Equal(t, count, 1, "")
}

package taskqueue_test

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/golibrary/stream/taskqueue"
)

const (
	parallelWorkSubmissions = 10000
	workTotal               = 49995000
)

var (
	prev    int
	counter int
	total   int32
	count   int32
)

func TestSerialQueue(t *testing.T) {
	q := taskqueue.New(taskqueue.Depth(100), taskqueue.Workers(1))
	prev = -1
	counter = 0
	for i := 0; i < 200; i++ {
		submitSerial(q, i)
	}
	q.Shutdown()
	assert.Equal(t, 199, prev)
	assert.Equal(t, 200, counter)
}

func submitSerial(q *taskqueue.Queue, i int) {
	q.Submit(func() {
		if i-1 == prev {
			prev = i
			counter++
		}
	})
}

func TestParallelQueue(t *testing.T) {
	q := taskqueue.New(taskqueue.Workers(5))
	total = 0
	count = 0
	for i := 0; i < parallelWorkSubmissions; i++ {
		submitParallel(q, i)
	}
	q.Shutdown()
	assert.Equal(t, parallelWorkSubmissions, int(count))
	assert.Equal(t, workTotal, int(total))
}

func submitParallel(q *taskqueue.Queue, i int) {
	q.Submit(func() {
		atomic.AddInt32(&total, int32(i))
		atomic.AddInt32(&count, 1)
	})
}

func boom() {
	var bad *int
	*bad = 1
}

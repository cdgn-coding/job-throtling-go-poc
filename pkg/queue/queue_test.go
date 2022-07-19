package queue

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type RunnableMock struct {
	mock.Mock
}

func (m *RunnableMock) run() {
	m.Called()
}

func TestNewQueue(t *testing.T) {
	assert.NotPanics(t, func() {
		NewQueue(5, 5)
	})
}

func TestQueueWorkers(t *testing.T) {
	t.Run("Queue workers with one job", func(t *testing.T) {
		queue := NewQueue(5, 5)
		queue.Start()

		done := make(chan bool)

		job := new(RunnableMock)
		job.On("run").Run(func(args mock.Arguments) {
			done <- true
		})

		queue.AddJob(job)
		assert.True(t, <-done)
	})

	t.Run("Queue workers with multiple jobs", func(t *testing.T) {
		numberOfJobs := 100
		queue := NewQueue(5, 5)
		queue.Start()

		job := new(RunnableMock)
		wg := sync.WaitGroup{}
		job.On("run").Run(func(args mock.Arguments) {
			wg.Done()
		})

		for i := 0; i < numberOfJobs; i++ {
			queue.AddJob(job)
			wg.Add(1)
		}

		wg.Wait()
		job.AssertNumberOfCalls(t, "run", numberOfJobs)
	})

	t.Run("Number of workers limit job concurrency", func(t *testing.T) {
		var numberOfJobs = 100
		var numberOfWorkers = 5
		var workingJobs int32 = 0
		var workingJobsHistory = make([]int32, 0)

		queue := NewQueue(numberOfWorkers, numberOfWorkers)
		queue.Start()

		job := new(RunnableMock)
		wg := sync.WaitGroup{}
		job.On("run").Run(func(args mock.Arguments) {
			atomic.AddInt32(&workingJobs, 1)
			time.Sleep(50 * time.Millisecond)
			workingJobsHistory = append(workingJobsHistory, workingJobs)
			atomic.AddInt32(&workingJobs, -1)
			wg.Done()
		})

		for i := 0; i < numberOfJobs; i++ {
			queue.AddJob(job)
			wg.Add(1)
		}

		wg.Wait()
		for _, measured := range workingJobsHistory {
			assert.LessOrEqual(t, int(measured), numberOfWorkers)
		}
	})
}

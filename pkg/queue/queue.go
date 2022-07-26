package queue

type IQueue interface {
	Start()
	AddJob(runnable Runnable)
}

type Queue struct {
	jobs    chan Runnable
	workers int
}

func NewQueue(size int, workers int) *Queue {
	return &Queue{
		jobs:    make(chan Runnable, size),
		workers: workers,
	}
}

func (q *Queue) Start() {
	for i := 0; i < q.workers; i++ {
		go q.worker()
	}
}

func (q *Queue) worker() {
	for {
		job := <-q.jobs
		job.Run()
	}
}

func (q *Queue) AddJob(runnable Runnable) {
	q.jobs <- runnable
}

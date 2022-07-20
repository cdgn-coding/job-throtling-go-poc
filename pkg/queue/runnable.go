package queue

type Runnable interface {
	run()
}

type RunnableFunc struct {
	f func()
}

func (r *RunnableFunc) run() {
	r.f()
}

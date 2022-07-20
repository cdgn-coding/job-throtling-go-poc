package queue

type Runnable interface {
	Run()
}

type RunnableFunc struct {
	f func()
}

func NewRunnable(f func()) RunnableFunc {
	return RunnableFunc{f}
}

func (r RunnableFunc) Run() {
	r.f()
}

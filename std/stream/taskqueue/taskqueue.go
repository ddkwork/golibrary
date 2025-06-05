package taskqueue

import (
	"runtime"
)

type Logger func(v ...any)

type Task func()

type Option func(*Queue)

type Queue struct {
	in      chan Task
	done    chan bool
	depth   int
	workers int
}

func Depth(depth int) Option {
	return func(q *Queue) { q.depth = depth }
}

func Workers(workers int) Option {
	return func(q *Queue) { q.workers = workers }
}

func New(options ...Option) *Queue {
	numCPU := runtime.NumCPU()
	q := &Queue{
		in:    make(chan Task, numCPU*2),
		done:  make(chan bool),
		depth: -1,
	}
	for _, option := range options {
		option(q)
	}
	if q.workers < 1 {
		q.workers = 1 + numCPU
	}
	go q.process()
	return q
}

func (q *Queue) Submit(task Task) {
	q.in <- task
}

func (q *Queue) Shutdown() {
	close(q.in)
	<-q.done
}

func (q *Queue) process() {
	var received, processed uint64

	var backlog []Task
	if q.depth > 0 {
		backlog = make([]Task, 0, q.depth)
	}

	ready := make(chan bool, q.workers)
	tasks := make(chan Task, q.workers)
	for range q.workers {
		go q.work(tasks, ready)
	}

outer:
	for {
	inner:
		select {
		case task := <-q.in:
			if task == nil {
				break outer
			}
			received++
			if len(backlog) == 0 {
				select {
				case tasks <- task:
					break inner
				default:
				}
			}
			if q.depth < 0 || len(backlog) < q.depth {
				backlog = append(backlog, task)
			} else {
				<-ready
				processed++
				tasks <- backlog[0]
				copy(backlog, backlog[1:])
				backlog[len(backlog)-1] = task
			}
		case <-ready:
			processed++
			if len(backlog) > 0 {
				tasks <- backlog[0]
				copy(backlog, backlog[1:])
				backlog[len(backlog)-1] = nil
				backlog = backlog[:len(backlog)-1]
			}
		}
	}

	for _, task := range backlog {
	drain:
		for {
			select {
			case tasks <- task:
				break drain
			case <-ready:
				processed++
			}
		}
	}
	for received != processed {
		<-ready
		processed++
	}
	close(tasks)
	q.done <- true
}

func (q *Queue) work(tasks <-chan Task, ready chan<- bool) {
	for task := range tasks {
		q.runTask(task)
		ready <- true
	}
}

func (q *Queue) runTask(task Task) {
	task()
}

package queue

type Task struct {
	ChatID int64
	Addr   string
}

type Queue struct {
	ch chan *Task
}

func NewQueue(buffer int) *Queue {
	return &Queue{
		make(chan *Task, buffer),
	}
}

func (q *Queue) Push(task *Task) {
	q.ch <- task
}

func (q *Queue) Channel() <-chan *Task {
	return q.ch
}

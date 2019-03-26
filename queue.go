package itertools

import "sync"

type Queue struct {
	sync.Mutex
	elements []interface{}
}

func (q *Queue) empty() bool {
	return len(q.elements) == 0
}

func (q *Queue) put(i interface{}) {
	q.Lock()
	q.elements = append(q.elements, i)
	q.Unlock()
}

func (q *Queue) get() interface{} {
	q.Lock()
	item := q.elements[0]
	q.elements = q.elements[1:]
	q.Unlock()
	return item
}

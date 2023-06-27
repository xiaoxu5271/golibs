package golibs

import "time"

type Queue struct {
	queue chan interface{}
}

func (q *Queue) Init(size int) *Queue {
	q.queue = make(chan interface{}, size)
	return q
}

func (q *Queue) EnQueue(itm interface{}, timeout time.Duration) bool {
	select {
	case q.queue <- itm:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (q *Queue) DeQueue() <-chan interface{} {
	return q.queue
}

func (q *Queue) DeQueueBlock() interface{} {
	itm := <-q.queue
	return itm
}

func (q *Queue) DeQueueTimed(timeout time.Duration) (interface{}, bool) {
	var itm interface{}
	select {
	case itm = <-q.queue:
		return itm, true
	case <-time.After(timeout):
		return nil, false
	}
}

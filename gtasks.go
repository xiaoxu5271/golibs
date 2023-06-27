package golibs

import (
	"reflect"
	"sync"
	"time"
)

type Tasks struct {
	rwlock   sync.RWMutex
	wg       sync.WaitGroup
	fun      interface{}
	types    []reflect.Type
	event    chan interface{}
	iEvent   int64
	eventMap map[int64]chan interface{}
	capacity chan uint8
	pending  chan []interface{}
	lock     sync.Mutex
	rtype    RunType
}

const (
	QUEUED_TASK RunType = iota
	UNQUEUED_TASK
)

type RunType int

func (t *Tasks) Init(fun interface{}, events int, capacity uint32, rtype RunType) *Tasks {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.event != nil {
		panic("Tasks can only be init once.")
	}

	rf := reflect.ValueOf(fun)
	if rf.Kind() != reflect.Func {
		panic("Argument should be a function")
	}
	t.fun = fun

	types := rf.Type()
	num := types.NumIn()
	num--
	if num < 0 {
		panic("Need at least on \"event Event\" argument.")
	}
	var event chan interface{}
	if reflect.TypeOf(event) != types.In(0) {
		panic("First argument must be \"chan interface{}\".")
	}

	t.types = make([]reflect.Type, 0, num)

	for i := 0; i < num; i++ {
		t.types = append(t.types, types.In(i+1))
	}
	t.event = make(chan interface{}, events)

	t.rtype = rtype
	switch t.rtype {
	case QUEUED_TASK:
		if capacity != 0 {
			t.capacity = make(chan uint8, 1)
			t.pending = make(chan []interface{}, capacity)
		} else {
			panic("capacity can not be zero in QUEUED_TASK model")
		}
	case UNQUEUED_TASK:
		t.eventMap = make(map[int64]chan interface{})
		if capacity != 0 {
			t.capacity = make(chan uint8, capacity)
		}
	}
	return t
}

func (t *Tasks) run(args ...interface{}) {
	var iEvt int64
	if len(t.types) != len(args) {
		panic("Argument count is not right!")
	}

	arguments := make([]reflect.Value, 0, len(args))
	if t.rtype == UNQUEUED_TASK {
		t.rwlock.Lock()
		t.iEvent++
		iEvt = t.iEvent
		t.eventMap[iEvt] = t.event
		arguments = append(arguments, reflect.ValueOf(t.eventMap[iEvt]))
		t.rwlock.Unlock()

	} else {
		arguments = append(arguments, reflect.ValueOf(t.event))
	}

	for i, v := range args {
		arguments = append(arguments, reflect.ValueOf(v))
		at := reflect.TypeOf(v)
		if at != t.types[i] {
			panic("Argument type is not right!")
		}
	}

	fun := reflect.ValueOf(t.fun)
	fun.Call(arguments)
	if t.rtype == UNQUEUED_TASK {
		t.rwlock.Lock()
		delete(t.eventMap, iEvt)
		t.rwlock.Unlock()
	}
}

func (t *Tasks) uncheckRun(args ...interface{}) {
	t.run(args...)
	if t.capacity != nil {
		<-t.capacity
	}
	t.wg.Done()
}

func (t *Tasks) goUncheckRun(args ...interface{}) {
	t.wg.Add(1)
	go t.uncheckRun(args...)
}

func (t *Tasks) Run(args ...interface{}) {
	if t.rtype == QUEUED_TASK {
		panic("can not us gTasks.Run in QUEUED_TASK model")
	}
	if t.capacity != nil {
		t.capacity <- 1
	}
	t.goUncheckRun(args...)
}

func (t *Tasks) TryRun(timeout time.Duration, args ...interface{}) bool {
	if t.rtype == QUEUED_TASK {
		panic("can not us gTasks.TryRun in QUEUED_TASK model")
	}
	if t.capacity != nil {
		select {
		case t.capacity <- 1:
			t.goUncheckRun(args...)
			return true
		case <-time.After(timeout):
			return false
		}
	} else {
		t.goUncheckRun(args...)
		return true
	}
}

func (t *Tasks) checkRun() {
	for {
		as := <-t.pending
		t.run(as...)

		t.lock.Lock()
		if len(t.pending) == 0 {
			if t.capacity != nil {
				<-t.capacity
			}
			t.wg.Done()
			t.lock.Unlock()
			return
		} else {
			t.lock.Unlock()
		}
	}
}

func (t *Tasks) QueuedRun(args ...interface{}) bool {
	if t.rtype == UNQUEUED_TASK {
		panic("can not us gTasks.QueuedRun in UNQUEUED_TASK model")
	}
	t.lock.Lock()
	defer t.lock.Unlock()

	var as []interface{}
	as = append(as, args...)

	select {
	case t.pending <- as:
		select {
		case t.capacity <- 1:
			t.wg.Add(1)
			go t.checkRun()
		default:
		}
		return true
	default:
		return false
	}
}

func (t *Tasks) Emit(event interface{}) {
	if t.rtype == UNQUEUED_TASK {
		t.rwlock.RLock()
		for _, chEvent := range t.eventMap {
			select {
			case chEvent <- event:
			default:
			}
		}
		t.rwlock.RUnlock()
	} else {
		select {
		case t.event <- event:
		default:
		}
	}
}

func (t *Tasks) EmitSync(event interface{}) {
	t.event <- event
}

func (t *Tasks) Join() {
	t.wg.Wait()
}

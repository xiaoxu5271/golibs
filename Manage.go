package golibs

import (
	"sync"
)

type Manage struct {
	tlock sync.RWMutex
	tasks []*Task
}

func (m *Manage) Init() *Manage {
	m.tasks = make([]*Task, 0)
	return m
}

func (m *Manage) RegistTask(task *Task) {
	m.tlock.Lock()
	defer m.tlock.Unlock()
	m.tasks = append(m.tasks, task)
}

func (m *Manage) Broadcast(event interface{}) {
	m.tlock.RLock()
	defer m.tlock.RUnlock()
	// Last to first.
	for i := len(m.tasks) - 1; i >= 0; i-- {
		t := m.tasks[i]
		t.Emit(event)
	}
}

func (m *Manage) CreateTask(fun interface{}, events int, capacity uint32, rtype RunType) *Task {
	t := new(Task)
	t.Tasks.Init(fun, events, capacity, rtype)
	t.Manage = m
	m.RegistTask(t)

	return t
}

func (m *Manage) Join() {
	m.tlock.RLock()
	defer m.tlock.RUnlock()
	for _, t := range m.tasks {
		t.Tasks.Join()
	}
}

type Task struct {
	*Manage
	Tasks
}

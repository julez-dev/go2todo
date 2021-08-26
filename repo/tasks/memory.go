package tasks

import (
	"context"
	"sync"
)

type InMemory struct {
	tasks map[string]*Task
	l     *sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		tasks: make(map[string]*Task),
		l:     &sync.RWMutex{},
	}
}

func (mem *InMemory) CreateTask(_ context.Context, task *Task) (*Task, error) {
	mem.l.Lock()
	defer mem.l.Unlock()

	mem.tasks[task.ID] = task

	return task, nil
}

func (mem *InMemory) UpdateTask(_ context.Context, task *Task) (*Task, error) {
	mem.l.Lock()
	defer mem.l.Unlock()

	for id := range mem.tasks {
		if id == task.ID {
			mem.tasks[id] = task
		}
	}

	return task, nil
}

func (mem *InMemory) GetTask(_ context.Context, search string) (*Task, error) {
	mem.l.RLock()
	defer mem.l.RUnlock()

	for id, list := range mem.tasks {
		if id == search {
			return list, nil
		}
	}

	return nil, ErrNotFound

}

func (mem *InMemory) GetTasks(_ context.Context, listID string) ([]*Task, error) {
	mem.l.RLock()
	defer mem.l.RUnlock()

	matchingTasks := []*Task{}
	for _, tasks := range mem.tasks {
		if listID == tasks.ListID {
			matchingTasks = append(matchingTasks, tasks)
		}
	}

	return matchingTasks, nil
}

func (mem *InMemory) GetAllTasks(_ context.Context) ([]*Task, error) {
	mem.l.RLock()
	defer mem.l.RUnlock()

	allTasks := []*Task{}

	for _, tasks := range mem.tasks {

		allTasks = append(allTasks, tasks)
	}

	return allTasks, nil

}

func (mem *InMemory) DeleteTask(_ context.Context, taskID string) error {
	mem.l.Lock()
	defer mem.l.Unlock()

	delete(mem.tasks, taskID)

	return nil
}

func (mem *InMemory) DeleteTasks(_ context.Context, listID string) error {
	mem.l.Lock()
	defer mem.l.Unlock()

	for _, task := range mem.tasks {
		if task.ListID == listID {
			delete(mem.tasks, task.ID)
		}
	}

	return nil
}

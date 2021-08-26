package lists

import (
	"context"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type InMemory struct {
	lists map[string]*List
	l     *sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		lists: make(map[string]*List),
		l:     &sync.RWMutex{},
	}
}

func ExampleInMemory() *InMemory {
	inMem := NewInMemory()

	inMem.l.Lock()
	defer inMem.l.Unlock()

	titles := [...]string{"My first task.", "Test 1", "Test 2"}

	for _, title := range titles {
		uuid := uuid.NewV4().String()
		inMem.lists[uuid] = &List{
			ID:   uuid,
			Name: title,
		}
	}

	return inMem
}

func (mem *InMemory) CreateList(_ context.Context, list *List) (*List, error) {
	mem.l.Lock()
	defer mem.l.Unlock()

	mem.lists[list.ID] = list

	return list, nil
}

func (mem *InMemory) UpdateList(_ context.Context, list *List) (*List, error) {
	mem.l.Lock()
	defer mem.l.Unlock()

	for id := range mem.lists {
		if id == list.ID {
			mem.lists[id] = list
		}
	}

	return list, nil
}

func (mem *InMemory) GetList(_ context.Context, search string) (*List, error) {
	mem.l.RLock()
	defer mem.l.RUnlock()

	for id, list := range mem.lists {
		if id == search {
			return list, nil
		}
	}

	return nil, ErrNotFound
}

func (mem *InMemory) GetLists(_ context.Context) ([]*List, error) {
	mem.l.RLock()
	defer mem.l.RUnlock()

	lists := make([]*List, 0, len(mem.lists))

	for _, list := range mem.lists {
		lists = append(lists, list)
	}

	return lists, nil
}

func (mem *InMemory) DeleteList(_ context.Context, id string) error {
	mem.l.Lock()
	defer mem.l.Unlock()

	delete(mem.lists, id)

	return nil
}

func (mem *InMemory) DeleteLists(context.Context) error {
	mem.l.Lock()
	defer mem.l.Unlock()

	for key := range mem.lists {
		delete(mem.lists, key)
	}

	return nil
}

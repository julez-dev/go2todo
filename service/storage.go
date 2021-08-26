package service

import (
	"context"
	"sort"
	"time"

	"github.com/julez-dev/go2todo/repo/lists"
	"github.com/julez-dev/go2todo/repo/tasks"
	uuid "github.com/satori/go.uuid"
)

type Interface interface {
}

type Storage struct {
	TasksRepo tasks.Interface
	ListsRepo lists.Interface
}

func NewStorage(tasksRepo tasks.Interface, listsRepo lists.Interface) *Storage {
	return &Storage{
		TasksRepo: tasksRepo,
		ListsRepo: listsRepo,
	}
}

func (s *Storage) StoreTask(ctx context.Context, task *tasks.Task) (*tasks.Task, error) {
	uuid := uuid.NewV4().String()
	task.ID = uuid
	task.CreatedAt = time.Now()

	return s.TasksRepo.CreateTask(ctx, task)
}

func (s *Storage) GetTask(ctx context.Context, taskID string) (*tasks.Task, error) {
	return s.TasksRepo.GetTask(ctx, taskID)
}

func (s *Storage) GetTasks(ctx context.Context, listID string) ([]*tasks.Task, error) {
	tasks, err := s.TasksRepo.GetTasks(ctx, listID)

	if err != nil {
		return nil, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (s *Storage) UpdateTask(ctx context.Context, task *tasks.Task) (*tasks.Task, error) {
	return s.TasksRepo.UpdateTask(ctx, task)
}

func (s *Storage) DeleteTask(ctx context.Context, taskID string) error {
	return s.TasksRepo.DeleteTask(ctx, taskID)
}

func (s *Storage) DeleteTasks(ctx context.Context, listID string) error {
	return s.TasksRepo.DeleteTasks(ctx, listID)
}

func (s *Storage) StoreList(ctx context.Context, list *lists.List) (*lists.List, error) {
	uuid := uuid.NewV4().String()
	list.ID = uuid
	list.CreatedAt = time.Now()

	return s.ListsRepo.CreateList(ctx, list)
}

func (s *Storage) GetList(ctx context.Context, listID string) (*lists.List, error) {
	return s.ListsRepo.GetList(ctx, listID)
}

func (s *Storage) GetLists(ctx context.Context) ([]*lists.List, error) {
	lists, err := s.ListsRepo.GetLists(ctx)

	if err != nil {
		return nil, err
	}

	sort.Slice(lists, func(i, j int) bool {
		return lists[i].CreatedAt.Before(lists[j].CreatedAt)
	})

	return lists, nil
}

func (s *Storage) UpdateList(ctx context.Context, list *lists.List) (*lists.List, error) {
	return s.ListsRepo.UpdateList(ctx, list)
}

func (s *Storage) DeleteList(ctx context.Context, listID string) error {
	err := s.TasksRepo.DeleteTasks(ctx, listID)

	if err != nil {
		return err
	}

	return s.ListsRepo.DeleteList(ctx, listID)
}

func (s *Storage) DeleteLists(ctx context.Context) error {
	lists, err := s.GetLists(ctx)

	if err != nil {
		return err
	}

	for _, list := range lists {
		err = s.DeleteTasks(ctx, list.ID)

		if err != nil {
			return err
		}
	}

	return nil
}

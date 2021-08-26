package tasks

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("task does not exist")

type Task struct {
	ID        string    `json:"id"`
	ListID    string    `json:"list_id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type Interface interface {
	CreateTask(context.Context, *Task) (*Task, error)
	UpdateTask(context.Context, *Task) (*Task, error)
	GetTask(context.Context, string) (*Task, error)
	GetTasks(context.Context, string) ([]*Task, error)
	GetAllTasks(context.Context) ([]*Task, error)
	DeleteTask(context.Context, string) error
	DeleteTasks(context.Context, string) error
}

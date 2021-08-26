package lists

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("list does not exist")

type List struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Interface interface {
	CreateList(context.Context, *List) (*List, error)
	UpdateList(context.Context, *List) (*List, error)
	GetList(context.Context, string) (*List, error)
	GetLists(context.Context) ([]*List, error)
	DeleteList(context.Context, string) error
	DeleteLists(context.Context) error
}

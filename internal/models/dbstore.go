package models

import (
	"context"

	pb "todo-app/api/v1/proto"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// Todo Data record
type TodoData struct {
	UserId      string              `json:"userId"`
	TaskId      string              `json:"taskId,omitempty"`
	Description string              `json:"description,omitempty"`
	State       pb.TodoState        `json:"state,omitempty"`
	DueDate     timestamp.Timestamp `json:"userID,omitempty"`
}

type TodoUser struct {
	Record map[string]TodoData `json:"record"`
}

// TodoDataStore contains methods to add, update, delete and get Todo records
type TodoUserStore interface {
	AddUser(ctx context.Context, userId string) error
	DeleteUser(ctx context.Context, userId string) error
	AddTodo(ctx context.Context, t *TodoData) (*TodoData, error)
	UpdateTodo(ctx context.Context, t *TodoData) (*TodoData, error)
	DeleteTodo(ctx context.Context, t *TodoData) error
	// List a specific TODO list
	ListTodo(ctx context.Context, userId, taskId string) (*TodoData, error)
	ListAllTodos(ctx context.Context, userId string) ([]*TodoData, error)
}

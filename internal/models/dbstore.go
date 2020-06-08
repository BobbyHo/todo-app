package models

import (
	"context"

	pb "todo-app/api/v1/proto"
)

// Todo Data record
type TodoData struct {
	UserId      string       `json:"userId,omitempty"`
	TaskId      string       `json:"taskId,omitempty"`
	Description string       `json:"description,omitempty"`
	State       pb.TodoState `json:"state,omitempty"`
	DueDate     string       `json:"userID,omitempty"`
}

type TodoUser struct {
	Record map[string]TodoData `json:"record"`
}

// TodoDataStore contains methods to add, update, delete and get Todo records
type TodoUserStore interface {
	// Add an user
	AddUser(ctx context.Context, userId string) error
	// Delete an user
	DeleteUser(ctx context.Context, userId string) error
	// check if an user exists in the DB
	FindUser(userId string) bool
	// Add a Todo task
	AddTodo(ctx context.Context, t *TodoData) (*TodoData, error)
	// Update aa Todo task
	UpdateTodo(ctx context.Context, t *TodoData) (*TodoData, error)
	// Delete a todo task
	DeleteTodo(ctx context.Context, t *TodoData) error
	// List a specific TODO list
	ListTodo(ctx context.Context, userId, taskId string) (*TodoData, error)
	// List All Todos
	ListAllTodos(ctx context.Context, userId string) ([]*TodoData, error)
}

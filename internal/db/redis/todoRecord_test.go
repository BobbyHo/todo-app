package redis

import (
	"context"
	"reflect"
	"testing"
	"todo-app/internal/models"
)

func TestAddUser(t *testing.T) {
	newTestClient := NewClient()

	err := newTestClient.Open()
	if err != nil {
		t.Error("Failed to connect to redis")
	}
	defer newTestClient.Close()

	err = newTestClient.TodoStore.AddUser(context.Background(), "test-user")
	if err != nil {
		t.Error("Failed to add user")
	}
}

func TestUpdateTodo(t *testing.T) {
	newTestClient := NewClient()

	err := newTestClient.Open()
	if err != nil {
		t.Error("Failed to connect to redis")
	}
	defer newTestClient.Close()

	testTodo := models.TodoData{
		UserId:      "test-user",
		TaskId:      "test-task-update",
		Description: "Test Update Todo",
	}

	_, err = newTestClient.TodoStore.AddTodo(context.Background(), &testTodo)
	if err != nil {
		t.Errorf("Failed to add todo err: %v", err.Error())
	}

	testTodo.Description = "Updated Test Todo Test"
	_, err = newTestClient.TodoStore.UpdateTodo(context.Background(), &testTodo)
	if err != nil {
		t.Errorf("Failed to update todo err: %v", err.Error())
	}

	result, err := newTestClient.TodoStore.ListTodo(context.Background(), testTodo.UserId, testTodo.TaskId)
	if err != nil {
		t.Errorf("Failed to update TODO task error: %v\n", err.Error())
	} else if !reflect.DeepEqual(&testTodo, result) {
		t.Errorf("Failed to update TODO task error: %v\n", err.Error())
	}

	t.Logf("Updated Todo Result: %v\n", *result)

	err = newTestClient.TodoStore.DeleteTodo(context.Background(), &testTodo)
	if err != nil {
		t.Errorf("Expected Result: %v is not the same as %v\n", testTodo, result)
	}
}

func TestAddandDeleteTodo(t *testing.T) {
	newTestClient := NewClient()

	err := newTestClient.Open()
	if err != nil {
		t.Error("Failed to connect to redis")
	}
	defer newTestClient.Close()

	testTodo := models.TodoData{
		UserId:      "test-user",
		TaskId:      "test-task-add-delete",
		Description: "Test Delete Todo",
	}
	_, err = newTestClient.TodoStore.AddTodo(context.Background(), &testTodo)
	if err != nil {
		t.Errorf("Failed to add todo err: %v", err.Error())
	}

	err = newTestClient.TodoStore.DeleteTodo(context.Background(), &testTodo)
	if err != nil {
		t.Errorf("Failed to delete todo err: %v", err.Error())
	}
}

func TestListAllTodos(t *testing.T) {

	newTestClient := NewClient()

	err := newTestClient.Open()
	if err != nil {
		t.Error("Failed to connect to redis")
	}
	defer newTestClient.Close()

	testUser := "test-user-lists"

	err = newTestClient.TodoStore.AddUser(context.Background(), testUser)
	if err != nil {
		t.Errorf("Failed to add user error: %v\n", err.Error())
	}

	testTodo1 := models.TodoData{
		UserId:      testUser,
		TaskId:      "test-task-list-all-1",
		Description: "List All Todos",
	}

	testTodo2 := models.TodoData{
		UserId:      testUser,
		TaskId:      "test-task-list-all-2",
		Description: "List All Todos 2",
	}

	_, err = newTestClient.TodoStore.AddTodo(context.Background(), &testTodo1)
	if err != nil {
		t.Errorf("Failed to add first todo err: %v", err.Error())
	}

	_, err = newTestClient.TodoStore.AddTodo(context.Background(), &testTodo2)
	if err != nil {
		t.Errorf("Failed to add second todo err: %v", err.Error())
	}

	results, err := newTestClient.TodoStore.ListAllTodos(context.Background(), testUser)
	if err != nil {
		t.Errorf("Failed to get all TODO lists: %v", err.Error())
	}

	if len(results) != 2 {
		t.Errorf("Invalid TODO list length: %v expected length: 2", len(results))
	}

	for i, r := range results {
		t.Logf("Todo %v: %v\n", i, *r)
		if !reflect.DeepEqual(testTodo1, *r) &&
			!reflect.DeepEqual(testTodo2, *r) {
			t.Errorf("Actual TODO: %v is not the same as expected Todo", *r)
		}
	}

	err = newTestClient.TodoStore.DeleteTodo(context.Background(), &testTodo1)
	if err != nil {
		t.Errorf("Failed to delete first todo err: %v", err.Error())
	}

	err = newTestClient.TodoStore.DeleteTodo(context.Background(), &testTodo2)
	if err != nil {
		t.Errorf("Failed to delete second todo err: %v", err.Error())
	}
}

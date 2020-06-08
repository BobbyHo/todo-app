package dbredis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	pb "todo-app/api/v1/proto"
	"todo-app/internal/models"

	redis "github.com/go-redis/redis"
)

// Ensure TodoUserStore implements models.TodoUserStore interface.
var _ models.TodoUserStore = &TodoUserStore{}

// TodoUserStore uses redis to store and retrieve user Todo records
type TodoUserStore struct {
	client *Client
}

// Check if an user exists in the DB
func (s *TodoUserStore) FindUser(userId string) bool {
	err := s.client.db.Get(userId).Err()
	if err != nil {
		// user record already created
		return false
	}

	return true
}

// AddUser creates a new User Todo Record in the database.
func (s *TodoUserStore) AddUser(ctx context.Context, userId string) error {

	var err error
	//check if the user record already exists
	err = s.client.db.Get(userId).Err()
	if err == nil {
		// user record already created
		return nil
	} else {
		log.Printf("db Get return: %v\n", err.Error())
	}

	// create an empty user record
	nr := models.TodoUser{}
	nr.Record = make(map[string]models.TodoData)

	data, err := json.Marshal(nr)
	if err != nil {
		log.Printf("Json Marhsal error: %v\n", err.Error())
		return err
	}

	err = s.client.db.Set(userId, data, 0).Err()
	if err != nil {
		log.Printf("DB set Error: %v\n", err.Error())
	}

	return err
}

func (s *TodoUserStore) DeleteUser(ctx context.Context, userId string) error {
	var err error
	//check if the user record already exists
	_, err = s.client.db.Get(userId).Bytes()
	if err != nil {
		// user record already created
		log.Printf("Failed to get user record: %v", err.Error())
		return err
	}

	err = s.client.db.Del(userId).Err()
	if err != nil {
		log.Printf("Failed to delete user record: %v", err.Error())
	}

	return err
}

func (s *TodoUserStore) AddTodo(ctx context.Context, t *models.TodoData) (*models.TodoData, error) {
	key := t.UserId
	err := s.client.db.Watch(func(tx *redis.Tx) error {
		val, err := tx.Get(key).Bytes()
		if err != nil && err != redis.Nil {
			return err
		}

		temp := &models.TodoUser{}
		if err := json.Unmarshal(val, temp); err != nil {
			log.Printf("Failed to Unmarshal record error: %v\n", err.Error())
			return err
		}

		// check if the user record already exists
		if _, ok := temp.Record[t.TaskId]; ok {
			return fmt.Errorf("Record %v already exists", t.TaskId)
		}

		temp.Record[t.TaskId] = *t

		data, err := json.Marshal(temp)
		if err != nil {
			log.Printf("Failed to Marshal record error: %v\n", err.Error())
			return err
		}

		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.Set(key, data, 0)
			return nil
		})
		return err
	}, key)

	return t, err
}

// Update a Todo Task
func (s *TodoUserStore) UpdateTodo(ctx context.Context, t *models.TodoData) (*models.TodoData, error) {
	key := t.UserId
	var tt models.TodoData
	err := s.client.db.Watch(func(tx *redis.Tx) error {
		val, err := tx.Get(key).Bytes()
		if err != nil && err != redis.Nil {
			fmt.Printf("Failed to get record %v\n", err.Error())
			return err
		}

		temp := &models.TodoUser{}
		if err := json.Unmarshal(val, temp); err != nil {
			log.Printf("Failed to Unmarshal record error: %v\n", err.Error())
			return err
		}

		log.Printf("todoRecord UPdateTodo UserId: %v\n", t.UserId)
		tt.UserId = t.UserId
		tt.TaskId = t.TaskId

		if t.Description == "" {
			tt.Description = temp.Record[t.TaskId].Description
		} else {
			tt.Description = t.Description
		}

		if t.DueDate == "" {
			tt.DueDate = temp.Record[t.TaskId].DueDate
		} else {
			tt.DueDate = t.DueDate
		}

		if t.State == pb.TodoState_UNDFINED {
			log.Printf("using original state: %v\n", temp.Record[t.TaskId].State)
			tt.State = temp.Record[t.TaskId].State
		} else {
			tt.State = t.State
		}

		temp.Record[t.TaskId] = tt

		data, err := json.Marshal(temp)
		if err != nil {
			log.Printf("Failed to Marshal record error: %v\n", err.Error())
			return err
		}

		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.Set(key, data, 0)
			return nil
		})
		return err
	}, key)

	if err != nil {
		fmt.Printf("Failed to update record %v\n", err.Error())
	}

	return &tt, err
}

// Delete a Todo Task
func (s *TodoUserStore) DeleteTodo(ctx context.Context, t *models.TodoData) error {
	key := t.UserId
	err := s.client.db.Watch(func(tx *redis.Tx) error {
		val, err := tx.Get(key).Bytes()
		if err != nil && err != redis.Nil {
			fmt.Printf("Failed to get record %v\n", err.Error())
			return err
		}

		temp := &models.TodoUser{}
		if err := json.Unmarshal(val, temp); err != nil {
			return err
		}

		if _, ok := temp.Record[t.TaskId]; ok {
			delete(temp.Record, t.TaskId)

			if _, ok := temp.Record[t.TaskId]; ok {
				fmt.Printf("Failed to remove record %v from user Todo list\n", t.TaskId)
				return fmt.Errorf("Failed to remove record %v from map", t.TaskId)
			} else {
				fmt.Printf("deleted record %v from user todo list\n", t.TaskId)
			}
		} else {
			return fmt.Errorf("Record %v does not exists", t.TaskId)
		}

		data, err := json.Marshal(temp)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.Set(key, data, 0)
			return nil
		})
		return err
	}, key)

	if err != nil {
		fmt.Printf("Failed to delete record %v\n", err.Error())
	}

	return err
}

// List a specific TODO task
func (s *TodoUserStore) ListTodo(ctx context.Context, userId, taskId string) (*models.TodoData, error) {

	//check if the user record already exists
	val, err := s.client.db.Get(userId).Bytes()
	if err != nil {
		// user record already created
		return &models.TodoData{}, err
	}

	temp := &models.TodoUser{}
	if err := json.Unmarshal(val, temp); err != nil {
		return &models.TodoData{}, err
	}

	result, ok := temp.Record[taskId]
	if !ok {
		return &models.TodoData{}, err
	}

	return &result, nil
}

// List All Todo Tasks belong to a user
func (s *TodoUserStore) ListAllTodos(ctx context.Context, userId string) ([]*models.TodoData, error) {
	var todoLists []*models.TodoData
	//check if the user record already exists
	val, err := s.client.db.Get(userId).Bytes()
	if err != nil {
		// user record already created
		return todoLists, err
	}

	temp := &models.TodoUser{}
	if err := json.Unmarshal(val, temp); err != nil {
		return todoLists, err
	}

	for _, r := range temp.Record {
		tempTodo := r
		//fmt.Printf("ListAllTodos: %v\n", tempTodo)
		todoLists = append(todoLists, &tempTodo)
	}

	return todoLists, nil
}

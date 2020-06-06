package main

import (
	"context"
	"log"
	"todo-app/internal/models"

	pb "todo-app/api/v1/proto"
)

// AddTodo implements TodoServer.AddTodo
func (s *server) AddTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {

	// TODO: authenticate user

	replyMsg := "Add Todo Task Successful"

	newTodo := &models.TodoData{
		UserId:      in.GetUserid(),
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State:       in.GetTodoBody().GetState(),
		DueDate:     *in.GetTodoBody().GetDueDate(),
	}

	_, err := Service.Store.TodoUserStore.AddTodo(ctx, newTodo)
	if err != nil {
		replyMsg = "Add Todo Task Fail"
		log.Printf("Failed to add Todo: %v\n", err.Error())
	}

	log.Printf("AddTodo:\n %v\n", newTodo)
	reply := &pb.TodoReply{}
	reply.Message = replyMsg
	reply.TodoBody = in.GetTodoBody()
	return reply, err
}

// Add a new user
func (s *server) AddUser(ctx context.Context, in *pb.UserRequest) (*pb.UserReply, error) {
	replyMsg := "Create User Successful"

	err := Service.Store.TodoUserStore.AddUser(ctx, in.GetUserid())
	if err != nil {
		replyMsg = "Add User Fail"
		log.Printf("Failed to add a new user: %v\n", err.Error())
	}

	log.Printf("Added new user: %v\n", in.GetUserid())
	reply := &pb.UserReply{}
	reply.Message = replyMsg
	return reply, err
}

// UpdateTodo implements TodoServer.UpdateTodo
func (s *server) UpdateTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	// TODO: authenticate user

	replyMsg := "Update Todo Task Successful"

	newTodo := &models.TodoData{
		UserId:      in.GetUserid(),
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State:       in.GetTodoBody().GetState(),
		DueDate:     *in.GetTodoBody().GetDueDate(),
	}

	_, err := Service.Store.TodoUserStore.UpdateTodo(ctx, newTodo)
	if err != nil {
		replyMsg = "Update Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	log.Printf("UpdateTodo:\n %v\n", newTodo)
	reply := &pb.TodoReply{}
	reply.Message = replyMsg
	reply.TodoBody = in.GetTodoBody()
	return reply, err
}

// DeleteTodo implements TodoServer.DeleteTodo
func (s *server) DeleteTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	// TODO: authenticate user
	replyMsg := "Delete Todo Task Successful"

	newTodo := &models.TodoData{
		UserId:      in.GetUserid(),
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State:       in.GetTodoBody().GetState(),
		DueDate:     *in.GetTodoBody().GetDueDate(),
	}

	err := Service.Store.TodoUserStore.DeleteTodo(ctx, newTodo)
	if err != nil {
		replyMsg = "Delete Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	log.Printf("UpdateTodo:\n %v\n", newTodo)
	reply := &pb.TodoReply{}
	reply.Message = replyMsg
	reply.TodoBody = in.GetTodoBody()
	return reply, err
}

// ListTodo implements TodoServer.ListTodo
func (s *server) ListTodo(ctx context.Context, in *pb.TodoListRequest) (*pb.TodoReply, error) {
	replyMsg := "List Todo Task Successful"

	// TODO authenticate user
	res, err := Service.Store.TodoUserStore.ListTodo(ctx, in.GetUserid(), in.GetTaskid())
	if err != nil {
		replyMsg = "List Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	reply := &pb.TodoReply{}
	reply.Message = replyMsg
	if err == nil {
		reply.TodoBody = &pb.TodoBody{
			Taskid:      res.TaskId,
			Description: res.Description,
			State:       res.State,
			DueDate:     &res.DueDate,
		}
	}

	return reply, nil
}

// List all Todos
func (s *server) ListAllTodos(ctx context.Context, in *pb.UserRequest) (*pb.TodoListAllReply, error) {
	replyMsg := "List Todo Task Successful"

	// TODO authenticate user

	res, err := Service.Store.TodoUserStore.ListAllTodos(ctx, in.GetUserid())
	if err != nil {
		replyMsg = "List Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	listAllReply := &pb.TodoListAllReply{}
	listAllReply.Message = replyMsg

	if err == nil {
		for _, r := range res {
			p := &pb.TodoBody{
				Taskid:      r.TaskId,
				Description: r.Description,
				State:       r.State,
				DueDate:     &r.DueDate,
			}
			listAllReply.Items = append(listAllReply.Items, p)
		}
	}

	return listAllReply, nil
}

// Listen to other TODO actions
func (s *server) ListenTodos(pb.Todo_ListenTodosServer) error {
	return nil
}

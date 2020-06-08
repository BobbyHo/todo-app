package main

import (
	"context"
	"fmt"
	"log"
	"todo-app/internal/models"

	pb "todo-app/api/v1/proto"
)

func authenticateUser(userId string) bool {
	//check if an user has been create
	return Service.Store.TodoUserStore.FindUser(userId)

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

// Add a new user
func (s *server) DeleteUser(ctx context.Context, in *pb.UserRequest) (*pb.UserReply, error) {
	replyMsg := "Delete User Successful"

	err := Service.Store.TodoUserStore.DeleteUser(ctx, in.GetUserid())
	if err != nil {
		replyMsg = "Delete User Fail"
		log.Printf("Failed to add a new user: %v\n", err.Error())
	}

	log.Printf("Delete user: %v\n", in.GetUserid())
	reply := &pb.UserReply{}
	reply.Message = replyMsg
	return reply, err
}

// AddTodo implements TodoServer.AddTodo
func (s *server) AddTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {

	if !authenticateUser(in.GetUserid()) {
		return &pb.TodoReply{}, fmt.Errorf("User %v does not exist\n", in.GetUserid())
	}

	replyMsg := "Add Todo Task Successful"

	newTodo := &models.TodoData{
		UserId:      in.GetUserid(),
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State:       in.GetTodoBody().GetState(),
		DueDate:     in.GetTodoBody().GetDueDate(),
	}

	pubMsg := newTodo.UserId + "add todo" + " task: " + newTodo.TaskId
	err := Service.MsgQ.Publish(ctx, models.TodoCmdChannel, pubMsg)
	if err != nil {
		log.Printf("Failed to publish msg: %v error: %v\n", pubMsg, err.Error())
	}

	_, err = Service.Store.TodoUserStore.AddTodo(ctx, newTodo)
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

// UpdateTodo implements TodoServer.UpdateTodo
func (s *server) UpdateTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	if !authenticateUser(in.GetUserid()) {
		return &pb.TodoReply{}, fmt.Errorf("User %v does not exist\n", in.GetUserid())
	}

	replyMsg := "Update Todo Task Successful"

	newTodo := &models.TodoData{
		UserId:      in.GetUserid(),
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State:       in.GetTodoBody().GetState(),
		DueDate:     in.GetTodoBody().GetDueDate(),
	}

	log.Printf("UpdateTodo: UserId: %v\n", in.GetUserid())

	pubMsg := newTodo.UserId + " update todo" + " task: " + newTodo.TaskId
	err := Service.MsgQ.Publish(ctx, models.TodoCmdChannel, pubMsg)
	if err != nil {
		log.Printf("Failed to publish msg: %v error: %v\n", pubMsg, err.Error())
	}

	res, err := Service.Store.TodoUserStore.UpdateTodo(ctx, newTodo)
	if err != nil {
		replyMsg = "Update Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	log.Printf("UpdateTodo:\n %v\n", newTodo)
	reply := &pb.TodoReply{}
	reply.Message = replyMsg
	reply.TodoBody = &pb.TodoBody{
		Taskid:      res.TaskId,
		Description: res.Description,
		DueDate:     res.DueDate,
		State:       res.State,
	}
	return reply, err
}

// DeleteTodo implements TodoServer.DeleteTodo
func (s *server) DeleteTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	if !authenticateUser(in.GetUserid()) {
		return &pb.TodoReply{}, fmt.Errorf("User %v does not exist\n", in.GetUserid())
	}

	replyMsg := "Delete Todo Task Successful"

	newTodo := &models.TodoData{
		UserId:      in.GetUserid(),
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State:       in.GetTodoBody().GetState(),
		DueDate:     in.GetTodoBody().GetDueDate(),
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
	if !authenticateUser(in.GetUserid()) {
		return &pb.TodoReply{}, fmt.Errorf("User %v does not exist\n", in.GetUserid())
	}

	replyMsg := "List Todo Task Successful"

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
			DueDate:     res.DueDate,
		}
	}

	return reply, nil
}

// List all Todos
func (s *server) ListAllTodos(ctx context.Context, in *pb.UserRequest) (*pb.TodoListAllReply, error) {
	if !authenticateUser(in.GetUserid()) {
		return &pb.TodoListAllReply{}, fmt.Errorf("User %v does not exist\n", in.GetUserid())
	}

	replyMsg := "List Todo Task Successful"

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
				DueDate:     r.DueDate,
			}
			listAllReply.Items = append(listAllReply.Items, p)
		}
	}

	return listAllReply, nil
}

// Listen to other TODO actions
func (s *server) ListenTodos(in *pb.ListenRequest, stream pb.Todo_ListenTodosServer) error {
	if !authenticateUser(in.GetUserid()) {
		return fmt.Errorf("User %v does not exist\n", in.GetUserid())
	}

	//log.Println("ListenTodos")
	// create a new msq client and connect to pubsub
	//msgQ := msgredis.NewMsgHandler()
	msgQ := Service.MsgQ

	log.Printf("ListenTodos: Subscribe to message : %v\n", models.TodoCmdChannel)
	p, err := msgQ.Subscribe(context.Background(), models.TodoCmdChannel)
	if err != nil {
		log.Printf("Failed to subscribe to channel error: %v\n", err.Error())
		return err
	}
	defer p.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		log.Printf("ListenTodos: Waiting for message from: %v\n", models.TodoCmdChannel)
		msg := p.Receive(ctx)
		log.Printf("ListenTodos: Received a message from: %v\n", models.TodoCmdChannel)

		if msg != nil {
			listenReply := &pb.ListenReply{}
			listenReply.Action = msg.Payload
			log.Printf("ListenTodos: Sending message: %v\n", msg.Payload)
			if err := stream.Send(listenReply); err != nil {
				return err
			}
		}
	}

	return nil
}

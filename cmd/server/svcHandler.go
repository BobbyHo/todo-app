package main

import "todo-app/internal/models"

// AddTodo implements TodoServer.AddTodo
func (s *server) AddTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {

	// TODO: authenticate user

	var replyMsg := "Add Todo Task Successful"

	newTodo := &models.TodoData{
		UserId:      in.GetUserId()
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State: 		 in.GetTodoBody().GetState(),
		DueDate:     in.GetTodoBody().GetDueDate(),
	}

	res, err := Service.Store.TodoUserStore.AddTodo(ctx, newTodo)
	if err != nil {
		replyMsg = "Add Todo Task Fail"
		log.Printf("Failed to add Todo: %v\n", err.Error())
	}

	log.Printf("AddTodo:\n %v\n", newTodo)
	reply = &pb.TodoReply{}
	reply.Message = replyMsg
	reply.TodoBody = in.GetTodoBody()
	return &reply, err
}

// UpdateTodo implements TodoServer.UpdateTodo
func (s *server) UpdateTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	// TODO: authenticate user

	var replyMsg := "Update Todo Task Successful"

	newTodo := &models.TodoData{
		UserId:      in.GetUserId()
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State: 		 in.GetTodoBody().GetState(),
		DueDate:     in.GetTodoBody().GetDueDate(),
	}

	res, err := Service.Store.TodoUserStore.UpdateTodo(ctx, newTodo)
	if err != nil {
		replyMsg = "Update Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	log.Printf("UpdateTodo:\n %v\n", newTodo)
	reply = &pb.TodoReply{}
	reply.Message = replyMsg
	reply.TodoBody = in.GetTodoBody()
	return &reply, err
}

// DeleteTodo implements TodoServer.DeleteTodo
func (s *server) DeleteTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	// TODO: authenticate user
	var replyMsg := "Delete Todo Task Successful"
	
	newTodo := &models.TodoData{
		UserId:      in.GetUserId()
		TaskId:      in.GetTodoBody().GetTaskid(),
		Description: in.GetTodoBody().GetDescription(),
		State: 		 in.GetTodoBody().GetState(),
		DueDate:     in.GetTodoBody().GetDueDate(),
	}

	err := Service.Store.TodoUserStore.DeleteTodo(ctx, newTodo)
	if err != nil {
		replyMsg = "Delete Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	log.Printf("UpdateTodo:\n %v\n", newTodo)
	reply = &pb.TodoReply{}
	reply.Message = replyMsg
	reply.TodoBody = in.GetTodoBody()
	return &reply, err
}

// ListTodo implements TodoServer.ListTodo
func (s *server) ListTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	var replyMsg := "List Todo Task Successful"

	// TODO authenticate user
	res, err := Service.Store.TodoUserStore.ListTodo(ctx, in.GetUserId(), in.GetTaskId())
	if err != nil {
		replyMsg = "List Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	reply = &pb.TodoReply{}
	reply.Message = replyMsg
	if err == nil {
		reply.TodoBody = &pb.TodoBody{
			TaskId: res.TaskId,
			Description: res.Description,
			State: res.State
			DueDate: &res.DueDate
		}
	}
	
	return &reply, nil
}

// List all Todos
func (s *server) ListAllTodos(context.Context, *UserRequest) (*pb.TodoListAllReply, error) {
	var replyMsg := "List Todo Task Successful"

	// TODO authenticate user

	res, err := Service.Store.TodoUserStore.ListAllTodos(ctx, in.GetUserId())
	if err != nil {
		replyMsg = "List Todo Task Fail"
		log.Printf("Failed to update Todo: %v\n", err.Error())
	}

	listAllReply := &pb.TodoListAllReply{}
	listAllReply.Message = replyMsg

	if err == nil {
		for _, r := range res {
			p := &pb.TodoBody{
				TaskId: r.TaskId,
				Description: r.Description,
				State: r.State
				DueDate: &r.DueDate	
			}
			listAllReply.Items = append(listAllReply.Items, p)
		} 
	}
	
	return listAllReply, nil
}

// Listen to other TODO actions
func (s *server)ListenTodos(pb.Todo_ListenTodosServer) error {
	return nil
}
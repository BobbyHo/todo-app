package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	pb "todo-app/api/v1/proto"

	"google.golang.org/grpc"
)

const (
	address = "localhost:12345"
)

type TodoConfiguration struct {
	UserId      string  `json:"userId"`
	TaskId      *string `json:"taskId,omitempty"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"duedate,omitempty"`
	Progress    *int32  `json:"duedate,omitempty"`
}

func convertReplyToTodoConfiguration(res *pb.TodoReply) *TodoConfiguration {
	tc := TodoConfiguration{}
	tc.TaskId = &res.TodoBody.Taskid
	tc.Description = &res.TodoBody.Description
	tc.DueDate = &res.TodoBody.DueDate
	return tc
}

// JsonPrettyString returns JSON string for the TodoConfiguration
func (t *TodoConfiguration) JsonPrettyString() string {
	prettyJSON, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		fmt.Errorf("Failed to generate json: %v\n", err.Error())
		return ""
	}
	//fmt.Printf("%s\n", string(prettyJSON))
	return string(prettyJSON)
}

func (t *TodoConfiguration) createTodoRequest() (*pb.TodoRequest, error) {
	todoRequest := pb.TodoRequest{}
	todoRequest.Userid = t.UserId

	todoBody := &pb.TodoBody{}

	if t.TaskId != nil {
		todoBody.Taskid = *t.TaskId
	}

	if t.Description != nil {
		todoBody.Description = *t.Description
	}
	if t.DueDate != nil {
		todoBody.DueDate = *t.DueDate
	}

	if t.Progress != nil {
		todoBody.State = pb.TodoState(*t.Progress)
	} else {
		todoBody.State = pb.TodoState_UNDFINED
	}

	todoRequest.TodoBody = todoBody

	return &todoRequest, nil
}

func createConnection() (pb.TodoClient, *grpc.ClientConn, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	c := pb.NewTodoClient(conn)

	return c, conn, err

}

// AddUser ...
func (t *TodoConfiguration) AddUser() error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTodoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// add a new user
	newUser := pb.UserRequest{}
	newUser.Userid = t.UserId
	_, err = c.AddUser(ctx, &newUser)
	if err != nil {
		log.Fatalf("could not add a new user: %v", err)
	}

	return err

}

// AddTodo ...
func (t *TodoConfiguration) AddTodo() error {

	todoRequest, err := t.createTodoRequest()
	if err != nil {
		return err
	}

	c, conn, err := createConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.AddTodo(ctx, todoRequest)
	if err != nil {
		log.Fatalf("could not add Todo: %v", err)
	}

	return err
}

// UpdateTodo ...
func (t *TodoConfiguration) UpdateTodo() error {

	todoRequest, err := t.createTodoRequest()
	if err != nil {
		return err
	}

	c, conn, err := createConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.UpdateTodo(ctx, todoRequest)
	if err != nil {
		log.Fatalf("could not add Todo: %v", err)
	}

	return err
}

// DeleteTodo ...
func (t *TodoConfiguration) DeleteTodo() error {

	todoRequest, err := t.createTodoRequest()
	if err != nil {
		return err
	}

	c, conn, err := createConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.DeleteTodo(ctx, todoRequest)
	if err != nil {
		log.Fatalf("could not add Todo: %v", err)
	}

	return err
}

// UpdateTodo ...
func (t *TodoConfiguration) ListTodo() (*TodoConTodoConfiguration, error) {

	todoRequest, err := t.createTodoRequest()
	if err != nil {
		return err
	}

	c, conn, err := createConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err = c.ListTodo(ctx, todoRequest)
	if err != nil {
		log.Fatalf("could not list Todo: %v", err)
	}

	tc := convertReplyToTodoConfiguration(res)
	return tc, err
}

func ListenTodos(ctx context.Context) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTodoClient(conn)

	userId := "test-new-user"

	newListen := pb.ListenRequest{}
	newListen.Userid = userId
	stream, err := c.ListenTodos(ctx, &newListen)
	if err != nil {
		log.Fatalf("%v.ListenTodos = _, %v", c, err)
	}

	for {
		listenReply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListenTodos(_) = _, %v", c, err)
		}
		log.Printf("Received Listen Reply: %v\n", listenReply)
	}

}

/*
func main() {
	ctxListen, cancelListen := context.WithCancel(context.Background())
	defer cancelListen()

	go ListenTodos(ctxListen)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTodoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userId := "test-new-user"

	// add a new user
	newUser := pb.UserRequest{}
	newUser.Userid = userId
	_, err = c.AddUser(ctx, &newUser)
	if err != nil {
		log.Fatalf("could not add a new user: %v", err)
	}

	todoRequest := pb.TodoRequest{}
	todoRequest.Userid = userId

	todoBody := &pb.TodoBody{}
	todoBody.Taskid = "Test"
	todoBody.Description = "Test Todo"
	todoBody.DueDate = &timestamp.Timestamp{}

	todoRequest.TodoBody = todoBody

	r, err := c.UpdateTodo(ctx, &todoRequest)
	if err != nil {
		log.Fatalf("could not add Todo: %v", err)
	}
	log.Printf("Todo Response: %s\n", r.GetMessage())

	time.Sleep(5 * time.Second)
}
*/

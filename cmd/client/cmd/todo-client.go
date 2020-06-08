package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "todo-app/api/v1/proto"

	"google.golang.org/grpc"
)

type TodoConfiguration struct {
	UserId      string  `json:"userId"`
	TaskId      *string `json:"taskId,omitempty"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"duedate,omitempty"`
	Progress    *int32  `json:"progress,omitempty"`
}

func convertReplyToTodoConfiguration(res *pb.TodoReply) *TodoConfiguration {
	tc := TodoConfiguration{}
	tc.TaskId = &res.TodoBody.Taskid
	tc.Description = &res.TodoBody.Description
	tc.DueDate = &res.TodoBody.DueDate
	progress := int32(res.TodoBody.State)
	tc.Progress = &progress
	return &tc
}

func convertReplyToAllTodoConfigurations(userId string, res *pb.TodoListAllReply) []*TodoConfiguration {
	tc := []*TodoConfiguration{}

	for _, r := range res.Items {
		todoConfig := &TodoConfiguration{}
		todoConfig.UserId = userId
		todoConfig.TaskId = &r.Taskid
		todoConfig.Description = &r.Description
		todoConfig.DueDate = &r.DueDate
		progress := int32(r.State)
		todoConfig.Progress = &progress

		tc = append(tc, todoConfig)
	}
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

func (t *TodoConfiguration) createTodoListRequest() *pb.TodoListRequest {
	tLR := pb.TodoListRequest{}
	tLR.Userid = t.UserId
	tLR.Taskid = *t.TaskId

	return &tLR

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
		log.Println("Setting state to UNDEFIND")
		todoBody.State = pb.TodoState_UNDFINED
	}

	todoRequest.TodoBody = todoBody

	return &todoRequest, nil
}

func createConnection() (pb.TodoClient, *grpc.ClientConn, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverIp, grpc.WithInsecure(), grpc.WithBlock())
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
	//conn, err := grpc.Dial(serverIp, grpc.WithInsecure(), grpc.WithBlock())
	log.Printf("Connecting to server %v\n", serverIp)
	conn, err := grpc.Dial(serverIp, grpc.WithInsecure())
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
func (t *TodoConfiguration) UpdateTodo() (*TodoConfiguration, error) {

	todoRequest, err := t.createTodoRequest()
	if err != nil {
		return &TodoConfiguration{}, err
	}

	c, conn, err := createConnection()
	if err != nil {
		return &TodoConfiguration{}, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.UpdateTodo(ctx, todoRequest)
	if err != nil {
		log.Printf("could not add Todo: %v", err.Error())
		return &TodoConfiguration{}, err
	}

	tc := convertReplyToTodoConfiguration(res)

	return tc, err
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

// ListAllTodos
func (t *TodoConfiguration) ListAllTodos() ([]*TodoConfiguration, error) {
	user := pb.UserRequest{}
	user.Userid = t.UserId

	reply := []*TodoConfiguration{}

	c, conn, err := createConnection()
	if err != nil {
		return reply, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.ListAllTodos(ctx, &user)
	if err != nil {
		log.Printf("could not list Todo: %v", err)
		return reply, err
	}

	tc := convertReplyToAllTodoConfigurations(t.UserId, res)
	return tc, err
}

// ListTodo ...
func (t *TodoConfiguration) ListTodo() (*TodoConfiguration, error) {

	todoRequest := t.createTodoListRequest()

	c, conn, err := createConnection()
	if err != nil {
		return &TodoConfiguration{}, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.ListTodo(ctx, todoRequest)
	if err != nil {
		log.Printf("could not list Todo: %v", err)
		return &TodoConfiguration{}, err
	}

	tc := convertReplyToTodoConfiguration(res)
	tc.UserId = t.UserId
	return tc, err
}

func (t *TodoConfiguration) ListenTodos() {
	// register for signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGHUP)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Set up a connection to the server.
	ctxListen, cancelListen := context.WithCancel(context.Background())
	defer cancelListen()

	conn, err := grpc.Dial(serverIp, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTodoClient(conn)

	newListen := pb.ListenRequest{}
	newListen.Userid = t.UserId
	stream, err := c.ListenTodos(ctxListen, &newListen)
	if err != nil {
		log.Fatalf("%v.ListenTodos = _, %v", c, err)
	}

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
	listenloop:
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Receive Done Signal in Listen")
				break listenloop
			default:
				listenReply, err := stream.Recv()
				if err == io.EOF {
					break listenloop
				}
				if err != nil {
					log.Fatalf("%v.ListenTodos(_) = _, %v", c, err)
				}
				fmt.Printf("Received Listen Reply: %v\n", listenReply)
			}
		}
	}(ctxListen)

waitloop:
	// wait for the stop signal
	for {
		select {
		case sig := <-sigs:
			log.Printf("\nReceived a signal: %v\n", sig)
			cancelListen() // signal the gRPC server to shutdown
			break waitloop
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	wg.Wait()

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

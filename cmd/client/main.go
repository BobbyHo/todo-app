/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"time"

	pb "todo-app/api/v1/proto"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

const (
	address = "localhost:12345"
)

func main() {
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

}

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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "todo-app/api/v1/proto"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

const (
	// TODO: make this configurable
	defaultPort = ":38888"
)

// server is used to implement todo.TodoServer
type server struct {
	pb.UnimplementedTodoServer
}

// server is used to implement todo.TodoServer.
type todo struct {
	userID      string              `json:"userId"`
	taskID      string              `json:"taskId,omitempty"`
	description string              `json:"description,omitempty"`
	dueDate     timestamp.Timestamp `json:"userID,omitempty"`
}

func (t *todo) String() string {
	s := "TODO Task: \n"

	s += fmt.Sprintf(" User ID: %s\n", t.userID)
	s += fmt.Sprintf(" Task ID: %s\n", t.taskID)
	s += fmt.Sprintf(" Description: %s\n", t.description)

	return s
}

// AddTodo implements todo.AddTodo
func (s *server) AddTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {

	newTodo := &todo{
		userID:      in.GetUserid(),
		taskID:      in.GetTaskid(),
		description: in.GetDescription(),
	}

	log.Printf("AddTodo:\n %v\n", newTodo)
	return &pb.TodoReply{Message: "AddTodo OK ", Requestbody: in}, nil
}

// UpdateTodo implements todo.AddTodo
func (s *server) UpdateTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	newTodo := &todo{
		userID:      in.GetUserid(),
		taskID:      in.GetTaskid(),
		description: in.GetDescription(),
	}

	log.Printf("UpdateTodo:\n %v\n", newTodo)
	return &pb.TodoReply{Message: "UpdateTodo OK ", Requestbody: in}, nil
}

// DeleteTodo implements todo.DeleteTodo
func (s *server) DeleteTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	newTodo := &todo{
		userID:      in.GetUserid(),
		taskID:      in.GetTaskid(),
		description: in.GetDescription(),
	}

	log.Printf("UpdateTodo:\n %v\n", newTodo)
	return &pb.TodoReply{Message: "UpdateTodo OK ", Requestbody: in}, nil
}

func main() {
	lis, err := net.Listen("tcp", defaultPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTodoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

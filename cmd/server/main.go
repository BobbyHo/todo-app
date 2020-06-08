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
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "todo-app/api/v1/proto"

	"github.com/golang/protobuf/ptypes/timestamp"
)

const (
	defaultPort = ":12345"
)

// server is used to implement todo.TodoServer
type server struct {
	pb.UnimplementedTodoServer
}

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

func main() {

	var configPath string

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	log.Printf("configuration Path: %v\n", configPath)

	// register for signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGHUP)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	Service = SvcInit(ctx, configPath)

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// no need to run SvcRun in a seperate go routine because
	// a go routine is created inside the handleTCPServer() function
	// defined in svc.go
	SvcRun(ctx, &wg, errc)

waitloop:
	// wait for the stop signal
	for {

		select {
		case sig := <-sigs:
			log.Printf("\nReceived a signal: %v\n", sig)
			cancel() // signal the gRPC server to shutdown
			Service.Database.Close()
			break waitloop
		case err := <-errc:
			log.Printf("Received error from gRPC server: %v", err.Error())
			cancel()
			//wg.Wait()
			os.Exit(0)
		}
	}

	wg.Wait()

	log.Print("Exiting Todo Service - Bye\n")

}

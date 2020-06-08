package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
	dbredis "todo-app/internal/db/redis"
	"todo-app/internal/models"
	msgredis "todo-app/internal/msgqueue/redis"

	pb "todo-app/api/v1/proto"

	"google.golang.org/grpc"
)

var (
	stop = make(chan bool)
	done = make(chan bool)
)

// GlobalConfig defines the configuration parameters for this web service
type GlobalConfig struct {
	Addr      string `json:"host,omitempty"`      //default value is empty string
	Port      int    `json:"port,omitempty"`      //default value is 12345
	DBAddress string `json:"dbaddress,omitempty"` //default value is 127.0.0.1
	DBPort    int    `json:"dbport,omitempty"`    //default value 6379 -- REDIS default port
}

var globalConfig GlobalConfig
var Service *models.Service

func readGlobalConf(conf string) (*GlobalConfig, error) {
	confData, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Printf("Failed to read config file (%v) err: %v\n", conf, err.Error())
		return nil, err
	}

	tempGlobalConf := GlobalConfig{}

	err = json.Unmarshal(confData, &tempGlobalConf)
	if err != nil {
		log.Printf("Failed to parse config File (%v) error: %v\n", conf, err.Error())
		return nil, err
	}

	return &tempGlobalConf, nil
}

func handleGrpcServer(ctx context.Context, address string, wg *sync.WaitGroup, errc chan error) {

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTodoServer(s, &server{})

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start GRPC server in a separate goroutine.
		go func() {
			log.Printf("GRPC server listening on %v", address)
			errc <- s.Serve(lis)
		}()

		<-ctx.Done()
		log.Printf("shutting down GRPC server at %q", address)

		stopped := make(chan struct{})
		go func() {
			s.GracefulStop()
			close(stopped)
		}()

		t := time.NewTimer(10 * time.Second)
		select {
		case <-t.C:
			s.Stop()
		case <-stopped:
			t.Stop()
		}

	}()
}

func openService(ctx context.Context, dbaddress, password string) *models.Service {
	db := dbredis.NewClient()

	if err := db.Open(dbaddress, password); err != nil {
		log.Fatalf("Error opening database: %v", err.Error())
		//panic("Error opening database")
	}

	// connect to a message queue to publish Todo events
	msgQ := msgredis.NewMsgHandler()

	// TODO: create a separate configuration for the message queue in the future
	if err := msgQ.Open(dbaddress, password); err != nil {
		log.Fatalf("Error opening message queue: %v", err.Error())
	}

	return &models.Service{
		Database: db,
		MsgQ:     msgQ,
		Store: &models.Store{
			TodoUserStore: db.TodoUserStore,
		},
	}
}

// SvcInit initializes the service configuration and DB connections
func SvcInit(ctx context.Context, configPath string) *models.Service {
	// set default values
	globalConfig.Addr = ""
	globalConfig.Port = 12345
	globalConfig.DBAddress = "127.0.0.1"
	globalConfig.DBPort = 6379

	if configPath != "" {
		tempConfig, err := readGlobalConf(configPath)
		if err == nil {
			globalConfig = *tempConfig
		}
	}

	dbaddress := globalConfig.DBAddress + ":" + strconv.Itoa(globalConfig.DBPort)
	// for now we assume that there is no password required
	return openService(ctx, dbaddress, "")
}

// SvcRun starts the TODO Service
func SvcRun(ctx context.Context, wg *sync.WaitGroup, errc chan error) {

	address := globalConfig.Addr + ":" + strconv.Itoa(globalConfig.Port)

	handleGrpcServer(ctx, address, wg, errc)
}

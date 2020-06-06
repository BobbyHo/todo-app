package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	stop = make(chan bool)
	done = make(chan bool)
)

// GlobalConfig defines the configuration parameters for this web service
type GlobalConfig struct {
	Addr      string `json:"host,omitempty"`       //default value is 0.0.0.0
	Port      int    `json:"port,omitempty"`       //default value is 12345
	DBAddress string `json:"db-address,omitempty"` //default value is 127.0.0.1
	DBPort    int    `json:"db-address,omitempty"` //default value 6379 -- REDIS default port
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
		log.Printf("shutting down GRPC server at %q", u.Host)

		stopped := make(chan, struct{})
		go func () {
			srv.GracefulStop()
			close(stopped)
		}

		t := time.NewTimer(10 * time.Second)
		select {
		case <- t.C:
			srv.Stop()
		case <- stopped:
			t.Stop()
		}

	}()
}

func openService(ctx context.Context, dbaddress, password string) *models.Service {
	db := redis.NewClient()

	if err := db.Open(dbaddress, password); err != nil {
		log.Fatalf("Error opening database: %v", err.Error())
		//panic("Error opening database")
	}

	return &models.Service{
		Database:      db,
		Store: &models.Store{
			TodoUserStore: db.TodoUserStore,
		},
	}
}

// SvcInit initializes the service configuration and DB connections
func SvcInit(ctx context.Context, configPath string) *models.Service {
	// set default values
	globalConfig.Addr = "0.0.0.0"
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

package redis

import (
	"fmt"
	"log"
	"time"

	redis "github.com/go-redis/redis"
)

const (
	// ErrUnableToOpen means we had an issue establishing a connection (or creating the database)
	ErrUnableToOpen = "Unable to open redis db; is there a core already running?  %v"

	// ErrUnableToInitialize means we couldn't initialize the redis DB
	ErrUnableToInitialize = "Unable to boot redis db:  %v"
)

// Client is a client for the redis data store.
type Client struct {
	db            *redis.Client
	Now           func() time.Time
	lastWriteTime time.Time
	TodoStore     *TodoUserStore
}

// NewClient initializes all stores
func NewClient() *Client {
	c := &Client{Now: time.Now}
	c.TodoStore = &TodoUserStore{client: c}
	return c
}

// Close the connection to the redis database
func (c *Client) Open() error {
	c.db = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", //todo: read this from env
		Password: "",
		DB:       0,
	})

	if c.db == nil {
		return fmt.Errorf("Failed to open redis connection")
	}

	return nil
}

// Close the connection to the redis database
func (c *Client) Close() error {
	if c.db != nil {
		log.Print("Closing Database connetion")
		return c.db.Close()
	}
	return nil
}

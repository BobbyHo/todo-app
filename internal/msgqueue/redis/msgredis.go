// using redis as a message broker for pub/sub
package msgredis

import (
	"context"
	"fmt"
	"log"
	"time"
	"todo-app/internal/models"

	"github.com/go-redis/redis"
)

// TodoMsgHandlers uses redis to publish and receive user Todo actions
type TodoMsgHandler struct {
	db            *redis.Client
	pubsub        *redis.PubSub
	Now           func() time.Time
	lastWriteTime time.Time
}

// NewClient initializes all stores
func NewMsgHandler() *TodoMsgHandler {
	c := &TodoMsgHandler{Now: time.Now}

	return c
}

// Close the connection to the redis database
func (c *TodoMsgHandler) Open(dbaddress, password string) error {
	c.db = redis.NewClient(&redis.Options{
		Addr:     dbaddress,
		Password: password,
		DB:       0,
	})

	if c.db == nil {
		return fmt.Errorf("Failed to open redis connection")
	}

	return nil
}

// Close the connection to the redis database
func (c *TodoMsgHandler) Close() error {
	if c.db != nil {
		log.Print("Closing Database connetion")
		return c.db.Close()
	}
	return nil
}

// Publish implements models.Publish
func (s *TodoMsgHandler) Publish(ctx context.Context, channel string, message interface{}) error {
	return s.db.Publish(channel, message).Err()
}

// Publish implements models.Publish
func (s *TodoMsgHandler) Subscribe(channels ...string) error {
	pubsub := s.db.Subscribe(channels...)
	s.pubsub = pubsub
	return nil
}

// Receive a message, it will block until a message has arrived
func (s *TodoMsgHandler) Receive(ctx context.Context) *models.Message {
	if s.pubsub == nil {
		return nil
	}

	/*
		msg, err := s.pubsub.ReceiveMessage()
		if err != nil {
			log.Printf("Failed to receive message: %v\n", err.Error())
		}
	*/

	msgChannel := s.pubsub.Channel()
	todoMsg := models.Message{}

outerloop:
	for {
		select {
		case <-ctx.Done():
			log.Print("Received cancel signal in Receive")
			break
		case msg := <-msgChannel:
			log.Printf("Received a message: %v\n", *msg)
			todoMsg.Channel = msg.Channel
			todoMsg.Payload = msg.Payload
			break outerloop
		}
	}

	return &todoMsg

}

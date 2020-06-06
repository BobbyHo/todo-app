package models

import "context"

// channels for the todo service
const (
	TodoCmdChannel = "Todo-CMD"
)

// TodoDataStore contains methods to add, update, delete and get Todo records
type TodoMsg interface {
	// Publish a TODO event
	Publish(ctx context.Context, channel string, message interface{}) error
	// Subscribe to a TODO event
	Subscribe(ctx context.Context, channel string, message interface{}) error
	//Receive A message from the subscribed channels
	Receive(ctx context.Context) interface{}
}

/*
// Pubsub Methods
type PubSub interface {
	//Close a pubsub handler
	Close() error
}
*/

// Todo Message
type Message struct {
	Channel string
	Payload string
}

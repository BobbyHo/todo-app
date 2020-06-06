package redis

import (
	"testing"
)

func TestClientOpen(t *testing.T) {

	address := "127.0.0.1:6379"
	newTestClient := NewClient()

	err := newTestClient.Open(address, "")
	if err != nil {
		t.Error("Failed to connect to redis")
	}
	newTestClient.Close()
}

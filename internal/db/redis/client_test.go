package redis

import (
	"testing"
)

func TestClientOpen(t *testing.T) {

	newTestClient := NewClient()

	err := newTestClient.Open()
	if err != nil {
		t.Error("Failed to connect to redis")
	}
	newTestClient.Close()
}

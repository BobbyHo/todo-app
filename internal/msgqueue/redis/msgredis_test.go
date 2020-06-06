package msgredis

import (
	"context"
	"testing"
	"todo-app/internal/models"
)

// Start Redis first before running the following tests
var testDBAddress string = "127.0.0.1:6379"
var testDBPassword string = ""

func TestPubSub(t *testing.T) {
	testClient := TodoMsgHandler{}

	err := testClient.Open(testDBAddress, testDBPassword)
	if err != nil {
		t.Error("Failed to connect to redis")
	}
	defer testClient.Close()

	/*
		testTodo := models.TodoData{
			UserId:      "test-user",
			TaskId:      "test-task-update",
			Description: "Test Update Todo",
		}
	*/

	err = testClient.Subscribe(models.TodoCmdChannel)
	if err != nil {
		t.Error("Invalid Pubsub Handler")
	}
	//defer testPubSub.Close()

	err = testClient.db.Publish(models.TodoCmdChannel, "hello").Err()
	if err != nil {
		t.Errorf("Failed to publish message: %v\n", err.Error())
	}

	testMSG := testClient.Receive(context.Background())

	t.Logf("Received message is: %v\n", testMSG)

}

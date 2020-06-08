package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configMap map[string]string = make(map[string]string)

// addtodoCmd represents the addtodo command
var addtodoCmd = &cobra.Command{
	Use:   "addtodo",
	Short: "add a todo task",
	Long:  `Add a new todo task`,
	Run: func(cmd *cobra.Command, args []string) {
		addTodo()
	},
}

func init() {
	rootCmd.AddCommand(addtodoCmd)
	addtodoCmd.MarkFlagRequired("username")
	addtodoCmd.MarkFlagRequired("task")
	addtodoCmd.MarkFlagRequired("serverip")
}

func addTodo() {
	if userId == "" {
		fmt.Println("Must provide user name in addtodo command")
		return
	}

	if taskId == "" {
		fmt.Println("Must provide task name in addtodo command")
		return
	}

	todoConfig := TodoConfiguration{}
	todoConfig.UserId = userId
	todoConfig.TaskId = &taskId

	if description != "-1-1" {
		todoConfig.Description = &description
	}

	if duedate != "-1-1" {
		todoConfig.DueDate = &duedate
	}

	if progress != -1 {
		todoConfig.Progress = &progress
	} else {
		var p int32 = 0
		todoConfig.Progress = &p
	}

	err := todoConfig.AddTodo()
	if err != nil {
		fmt.Printf("Failed to add todo: error: %v\n", err.Error())
	} else {
		jsonStr := todoConfig.JsonPrettyString()
		fmt.Printf("Added todo successfully:\n%v\n", jsonStr)
	}
}

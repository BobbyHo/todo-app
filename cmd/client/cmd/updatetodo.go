package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// updatetodoCmd represents the updatetodo command
var updatetodoCmd = &cobra.Command{
	Use:   "updatetodo",
	Short: "update an existing todo task",
	Long:  `Update an existing Todo task`,
	Run: func(cmd *cobra.Command, args []string) {
		updateTodo()
	},
}

func init() {
	rootCmd.AddCommand(updatetodoCmd)
	updatetodoCmd.MarkFlagRequired("username")
	updatetodoCmd.MarkFlagRequired("task")
}

func updateTodo() {
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
	}

	tc, err := todoConfig.UpdateTodo()
	if err != nil {
		fmt.Printf("Failed to update todo: error: %v\n", err.Error())
	} else {
		tc.UserId = todoConfig.UserId
		jsonStr := tc.JsonPrettyString()
		fmt.Printf("Updated todo successfully: \n %v\n", jsonStr)
	}
}

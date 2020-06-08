package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listtodoCmd represents the listtodo command
var listtodoCmd = &cobra.Command{
	Use:   "listtodo",
	Short: "list a todo task",
	Long:  `List a todo task`,
	Run: func(cmd *cobra.Command, args []string) {
		listTodo()
	},
}

func init() {
	rootCmd.AddCommand(listtodoCmd)
	addtodoCmd.MarkFlagRequired("username")
	addtodoCmd.MarkFlagRequired("task")
}

func listTodo() {
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

	tc, err := todoConfig.ListTodo()
	if err != nil {
		fmt.Printf("Failed to list todo: error: %v\n", err.Error())
	} else {
		jsonStr := tc.JsonPrettyString()
		fmt.Printf("%s\n", jsonStr)
	}
}

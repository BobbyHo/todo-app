package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listtodoCmd represents the listtodo command
var listtodoCmd = &cobra.Command{
	Use:   "listtodo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

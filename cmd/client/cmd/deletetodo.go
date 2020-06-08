package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deletetodoCmd represents the deletetodo command
var deletetodoCmd = &cobra.Command{
	Use:   "deletetodo",
	Short: "delete a todo task",
	Long:  `Delete an existing Todo Task`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteTodo()
	},
}

func init() {
	rootCmd.AddCommand(deletetodoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deletetodoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deletetodoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteTodo() {
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

	err := todoConfig.DeleteTodo()
	if err != nil {
		fmt.Printf("Failed to delete todo: error: %v\n", err.Error())
	} else {
		jsonStr := todoConfig.JsonPrettyString()
		fmt.Printf("Deleted todo successfully: \n %v\n", jsonStr)
	}
}

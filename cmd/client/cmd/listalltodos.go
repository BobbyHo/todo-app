package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listalltodosCmd represents the listalltodos command
var listalltodosCmd = &cobra.Command{
	Use:   "listalltodos",
	Short: "list all todos",
	Long:  `List todo tasks belong to a user`,
	Run: func(cmd *cobra.Command, args []string) {
		listallTodo()
	},
}

func init() {
	rootCmd.AddCommand(listalltodosCmd)
	addtodoCmd.MarkFlagRequired("username")
}

func listallTodo() {
	if userId == "" {
		fmt.Println("Must provide user name in addtodo command")
		return
	}

	todoConfig := TodoConfiguration{}
	todoConfig.UserId = userId

	tc, err := todoConfig.ListAllTodos()
	if err != nil {
		fmt.Printf("Failed to list all todo: error: %v\n", err.Error())
	} else {
		for _, r := range tc {
			jsonStr := r.JsonPrettyString()
			fmt.Printf("%s\n", jsonStr)
		}
	}
}

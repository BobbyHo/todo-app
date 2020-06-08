package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listentodoCmd represents the listentodo command
var listentodoCmd = &cobra.Command{
	Use:   "listentodo",
	Short: "listen to todo commands",
	Long:  `Listen to todo commands`,
	Run: func(cmd *cobra.Command, args []string) {
		listenTodos()
	},
}

func init() {
	rootCmd.AddCommand(listentodoCmd)
	listentodoCmd.MarkFlagRequired("username")
}

func listenTodos() {
	if userId == "" {
		fmt.Println("Must provide user name in addtodo command")
		return
	}

	todoConfig := TodoConfiguration{}
	todoConfig.UserId = userId

	todoConfig.ListenTodos()
}

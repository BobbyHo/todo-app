package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// adduserCmd represents the adduser command
var adduserCmd = &cobra.Command{
	Use:   "adduser",
	Short: "add a user to Todo App",
	Long:  `Register a new user to a Todo App - This is required before calling other Todo operations`,
	Run: func(cmd *cobra.Command, args []string) {
		addUser()
	},
}

func init() {
	rootCmd.AddCommand(adduserCmd)
	adduserCmd.MarkFlagRequired("username")
	adduserCmd.MarkFlagRequired("serverip")
}

func addUser() {

	if userId == "" {
		fmt.Println("Must provide user name in adduser command")
		return
	}
	todoConfig := TodoConfiguration{}
	todoConfig.UserId = userId

	err := todoConfig.AddUser()
	if err != nil {
		fmt.Printf("Failed to add user: error: %v\n", err.Error())
	} else {
		jsonStr := todoConfig.JsonPrettyString()
		fmt.Printf("Added user successfully:\n%v\n", jsonStr)
	}

}

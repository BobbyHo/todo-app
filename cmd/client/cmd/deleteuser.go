package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteuserCmd represents the deleteuser command
var deleteuserCmd = &cobra.Command{
	Use:   "deleteuser",
	Short: "delete a user",
	Long:  `Delete an existing user`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteUser()
	},
}

func init() {
	rootCmd.AddCommand(deleteuserCmd)
	adduserCmd.MarkFlagRequired("username")
	adduserCmd.MarkFlagRequired("serverip")

}

func deleteUser() {

	if userId == "" {
		fmt.Println("Must provide user name in adduser command")
		return
	}
	todoConfig := TodoConfiguration{}
	todoConfig.UserId = userId

	err := todoConfig.DeleteUser()
	if err != nil {
		fmt.Printf("Failed to delete user: error: %v\n", err.Error())
	} else {
		jsonStr := todoConfig.JsonPrettyString()
		fmt.Printf("Deleted user successfully: \n %v\n", jsonStr)
	}
}

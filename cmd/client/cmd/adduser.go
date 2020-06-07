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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adduserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adduserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addUser() {
	//userId := ""

	/*
		for i, n := range args {
			if i == 0 {
				userId = userId + n
			} else {
				userId = userId + " " + n
			}
		}
	*/
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
		fmt.Printf("Added user successfully: \n %v\n", jsonStr)
	}

	//fmt.Printf("addUser: %v\n", userId)
}

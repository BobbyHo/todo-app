/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
		if len(args) == 0 {
			fmt.Println("user name can not be empty")
		} else {
			addUser(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(adduserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adduserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adduserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addUser(args []string) {
	userId := ""

	for i, n := range args {
		if i == 0 {
			userId = userId + n
		} else {
			userId = userId + " " + n
		}
	}
	fmt.Printf("addUser: %v\n", userId)
}

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
	"log"

	"github.com/spf13/cobra"
)

/*
var userId string
var taskId string
var description string
var duedate string
var progress int32
*/

var configMap map[string]string = make(map[string]string)

// addtodoCmd represents the addtodo command
var addtodoCmd = &cobra.Command{
	Use:   "addtodo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("addtodo called")

		if cmd.Flags().Lookup("description") != nil {
			log.Printf("description flag is set")
		}

		addTodo()
	},
}

func init() {
	rootCmd.AddCommand(addtodoCmd)
	//addtodoCmd.Flags().StringVarP(&userId, "username", "u", "", "User Name")
	addtodoCmd.MarkFlagRequired("username")
	//addtodoCmd.Flags().StringVarP(&taskId, "task", "t", "", "Task Title")
	addtodoCmd.MarkFlagRequired("task")

	//addtodoCmd.Flags().StringVarP(&description, "description", "d", "", "Task Description")
	//addtodoCmd.Flags().StringVarP(&duedate, "duedate", "e", "", "Task Due Date")
	//addtodoCmd.Flags().Int32VarP(&progress, "progress", "p", 0, "Task Progress 0: Todo; 1: In Progress; 2: Done")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addtodoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addtodoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addTodo() {
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

	err := todoConfig.AddTodo()
	if err != nil {
		fmt.Printf("Failed to add todo: error: %v\n", err.Error())
	} else {
		jsonStr := todoConfig.JsonPrettyString()
		fmt.Printf("Added todo successfully: \n %v\n", jsonStr)
	}

	//fmt.Printf("addtodoCmd username: %s\n", userId)
	//fmt.Printf("addtodoCmd taskname: %s\n", taskId)

}

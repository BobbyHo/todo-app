package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var adduser string
var deleteuser string

var userId string
var taskId string
var description string
var duedate string
var progress int32

var serverIp string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todoclient",
	Short: "A simple Todo application",
	Long:  `This Todo application allows a user to register or delete a user, add/update/delete/list Todo tasks`,

	Run: func(cmd *cobra.Command, args []string) { fmt.Println("Hello TodoClient") },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is todoclient.yaml)")

	rootCmd.PersistentFlags().StringVarP(&userId, "username", "u", "", "User Name")
	rootCmd.PersistentFlags().StringVarP(&taskId, "task", "t", "", "Task Title")

	rootCmd.PersistentFlags().StringVarP(&description, "description", "d", "-1-1", "Task Description")
	rootCmd.PersistentFlags().StringVarP(&duedate, "duedate", "e", "-1-1", "Task Due Date")
	rootCmd.PersistentFlags().Int32VarP(&progress, "progress", "p", -1, "Task Progress 0: Todo; 1: In Progress; 2: Done")
	rootCmd.PersistentFlags().StringVarP(&serverIp, "serverip", "s", "127.0.0.1:12345", "Server Address")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in the current directory with name "todoclient.yaml"
		viper.AddConfigPath("./")
		viper.SetConfigName("todoclient.yaml")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

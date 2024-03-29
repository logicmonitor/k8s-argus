package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var kubeConfigFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{ // nolint: exhaustivestruct
	Use:   "argus",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err) // nolint: forbidigo
		os.Exit(-1)
	}
}

// nolint: gochecknoinits
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&kubeConfigFile, "kubeconfig", "", "kubeconfig file path")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

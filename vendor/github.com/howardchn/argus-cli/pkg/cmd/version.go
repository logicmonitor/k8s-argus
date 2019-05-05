package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "argus-cli version",
	Long:  "argus-cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("argus-cli", "v0.1.0-alpha")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

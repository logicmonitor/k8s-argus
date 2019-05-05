package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	accessId  string
	accessKey string
)

var RootCmd = cobra.Command{
	Use:   "argus-cli [Command] [Options]",
	Short: "argus associate utility",
	Long:  "argus associate utility",
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&accessId, "accessId", "i", "", "access id that is generated from santaba")
	RootCmd.PersistentFlags().StringVarP(&accessKey, "accessKey", "k", "", "access key that is generated from santaba")
	RootCmd.MarkFlagRequired("accessId")
	RootCmd.MarkFlagRequired("accessKey")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

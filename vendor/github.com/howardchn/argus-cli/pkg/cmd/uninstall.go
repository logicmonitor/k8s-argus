package cmd

import (
	"fmt"
	"github.com/howardchn/argus-cli/pkg"
	"github.com/howardchn/argus-cli/pkg/conf"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var (
	cluster  string
	account  string
	parentId int32
	mode     string
	confFile string
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall argus related resources",
	Long:  "uninstall argus related resources in santaba and k8s",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := getConfiguration()
		if err != nil {
			log.Fatal(err)
		}

		client := uninstaller.NewClient(conf)
		err = client.Clean()
		if err != nil {
			fmt.Printf("uninstall failed. err = %v\n", err)
			return
		}

		fmt.Println("uninstall success")
	},
}

func init() {
	uninstallCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "cluster name")
	uninstallCmd.Flags().StringVarP(&account, "account", "a", "", "account name")
	uninstallCmd.Flags().Int32VarP(&parentId, "parentId", "g", 1, "parent group id, default: 1")
	uninstallCmd.Flags().StringVarP(&mode, "mode", "m", "all", "uninstall mode: [rest|helm|all], default: all")
	uninstallCmd.Flags().StringVarP(&confFile, "confFile", "f", "", "configure file (*.yaml)")
	RootCmd.AddCommand(uninstallCmd)
}

func getConfiguration() (*conf.LMConf, error) {
	if confFile != "" {
		_, err := os.Stat(confFile)
		if err != nil {
			return nil, err
		}

		confBuffer, err := ioutil.ReadFile(confFile)
		var conf conf.LMConf
		err = yaml.Unmarshal(confBuffer, &conf)
		if err != nil {
			return nil, err
		}

		return &conf, nil
	} else {
		conf := &conf.LMConf{AccessId: accessId, AccessKey: accessKey, Account: account, Cluster: cluster, ParentId: parentId, Mode: mode}
		return conf, nil
	}
}

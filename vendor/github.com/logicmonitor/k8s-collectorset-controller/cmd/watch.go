// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"

	"github.com/logicmonitor/k8s-collectorset-controller/pkg/config"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/controller"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/server"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/storage/inmem"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the application configuration.
		collectorsetconfig, err := config.GetConfig()
		if err != nil {
			log.Fatalf("Failed to get config: %v", err)
		}

		// TODO: storage.Storage should define a Chan() func.
		countChan := make(chan int, 1)
		// Instantiate the storage backend.
		storage := inmem.New(countChan)

		// Instantiate the CollectorSet controller.
		collectorsetcontroller, err := controller.New(collectorsetconfig, storage)
		if err != nil {
			log.Fatalf("Failed to create CollectorSet controller: %v", err)
		}

		// Create the CRD if it does not already exist.
		_, err = collectorsetcontroller.CreateCustomResourceDefinition()
		if err != nil && !apierrors.IsAlreadyExists(err) {
			log.Fatalf("Failed to create CRD: %v", err)
		}

		// Start the CollectorSet controller.
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		go collectorsetcontroller.Run(ctx) // nolint: errcheck

		// Start the gRPC server.
		srv := server.New(storage, countChan)
		go srv.Run()

		select {}
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

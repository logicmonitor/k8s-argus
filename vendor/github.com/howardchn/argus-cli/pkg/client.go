package uninstaller

import (
	"github.com/howardchn/argus-cli/pkg/conf"
	"github.com/howardchn/argus-cli/pkg/helm"
	"github.com/howardchn/argus-cli/pkg/rest"
	"log"
	"strings"
)

type Client struct {
	Conf       *conf.LMConf
	RestClient *rest.Client
	HelmClient *helm.Client
}

func NewClient(conf *conf.LMConf) *Client {
	return &Client{
		conf,
		rest.NewClient(conf),
		helm.NewClient(conf),
	}
}

func (client *Client) Clean() error {
	mode := strings.ToLower(client.Conf.Mode)

	var err error
	if mode == "all" || mode == "rest" {
		err = client.RestClient.Clean()
		if err != nil {
			log.Println("-- LM uninstall failed --", err)
			return err
		} else {
			log.Println("-- LM uninstall success --")
		}
	}

	if mode == "all" || mode == "helm" {
		err = client.HelmClient.Clean()
		if err != nil {
			log.Println("-- helm uninstall failed --", err)
			return err
		} else {
			log.Println("-- helm uninstall success --")
		}
	}

	return nil
}

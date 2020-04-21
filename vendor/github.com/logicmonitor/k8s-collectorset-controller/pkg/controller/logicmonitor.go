package controller

import (
	"net/http"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/config"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	log "github.com/sirupsen/logrus"
)

func newLMClient(collectorsetconfig *config.Config) (*client.LMSdkGo, error) {
	config := client.NewConfig()
	config.SetAccessID(&collectorsetconfig.ID)
	config.SetAccessKey(&collectorsetconfig.Key)
	domain := collectorsetconfig.Account + ".logicmonitor.com"
	config.SetAccountDomain(&domain)
	if collectorsetconfig.ProxyURL == "" {
		return client.New(config), nil
	}
	return newClientWithProxy(config, collectorsetconfig)
}

func newClientWithProxy(config *client.Config, collectorsetconfig *config.Config) (*client.LMSdkGo, error) {
	proxyURL, err := url.Parse(collectorsetconfig.ProxyURL)
	if collectorsetconfig.ProxyUser != "" {
		if collectorsetconfig.ProxyPass != "" {
			proxyURL.User = url.UserPassword(collectorsetconfig.ProxyUser, collectorsetconfig.ProxyPass)
		} else {
			proxyURL.User = url.User(collectorsetconfig.ProxyUser)
		}
	}

	if err != nil {
		return nil, err
	}
	log.Infof("Using http/s proxy: %s", collectorsetconfig.ProxyURL)
	httpClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	transport := httptransport.NewWithClient(config.TransportCfg.Host, config.TransportCfg.BasePath, config.TransportCfg.Schemes, &httpClient)
	authInfo := client.LMv1Auth(*config.AccessID, *config.AccessKey)
	cli := new(client.LMSdkGo)
	cli.Transport = transport
	cli.LM = lm.New(transport, strfmt.Default, authInfo)
	return cli, nil
}

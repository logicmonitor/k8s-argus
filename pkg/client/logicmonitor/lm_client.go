package logicmonitor

import (
	"net/http"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
)

func NewLMClient(argusConfig *config.Config) (*client.LMSdkGo, error) {
	conf := client.NewConfig()
	conf.SetAccessID(&argusConfig.ID)
	conf.SetAccessKey(&argusConfig.Key)
	domain := argusConfig.Account + ".logicmonitor.com"
	conf.SetAccountDomain(&domain)
	// conf.UserAgent = constants.UserAgentBase + constants.Version
	if argusConfig.ProxyURL == "" {
		if argusConfig.IgnoreSSL {
			return newLMClientWithoutSSL(conf)
		}

		return client.New(conf), nil
	}

	return newLMClientWithProxy(conf, argusConfig)
}

func newLMClientWithProxy(config *client.Config, argusConfig *config.Config) (*client.LMSdkGo, error) {
	proxyURL, err := url.Parse(argusConfig.ProxyURL)
	if err != nil {
		return nil, err
	}
	if argusConfig.ProxyUser != "" {
		if argusConfig.ProxyPass != "" {
			proxyURL.User = url.UserPassword(argusConfig.ProxyUser, argusConfig.ProxyPass)
		} else {
			proxyURL.User = url.User(argusConfig.ProxyUser)
		}
	}
	httpClient := http.Client{
		Transport: &http.Transport{ // nolint: exhaustivestruct
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	transport := httptransport.NewWithClient(config.TransportCfg.Host, config.TransportCfg.BasePath, config.TransportCfg.Schemes, &httpClient)
	authInfo := client.LMv1Auth(*config.AccessID, *config.AccessKey)
	clientObj := new(client.LMSdkGo)
	clientObj.Transport = transport
	clientObj.LM = lm.New(transport, strfmt.Default, authInfo)

	return clientObj, nil
}

func newLMClientWithoutSSL(config *client.Config) (*client.LMSdkGo, error) {
	opts := httptransport.TLSClientOptions{InsecureSkipVerify: true}
	httpClient, err := httptransport.TLSClient(opts)
	if err != nil {
		return nil, err
	}
	transport := httptransport.NewWithClient(config.TransportCfg.Host, config.TransportCfg.BasePath, config.TransportCfg.Schemes, httpClient)
	authInfo := client.LMv1Auth(*config.AccessID, *config.AccessKey)
	cli := new(client.LMSdkGo)
	cli.Transport = transport
	cli.LM = lm.New(transport, strfmt.Default, authInfo)

	return cli, nil
}

package controller

import lm "github.com/logicmonitor/lm-sdk-go"

func newLMClient(id, key, company string) *lm.DefaultApi {
	config := lm.NewConfiguration()
	config.APIKey = map[string]map[string]string{
		"Authorization": {
			"AccessID":  id,
			"AccessKey": key,
		},
	}
	config.BasePath = "https://" + company + ".logicmonitor.com/santaba/rest"

	api := lm.NewDefaultApi()
	api.Configuration = config

	return api
}

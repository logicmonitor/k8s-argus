# lm-sdk-go
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/logicmonitor/lm-sdk-go)

Getting Started
---------------

```go
package main

import lmv1 "github.com/logicmonitor/lm-sdk-go"

func NewLMClient(id, key, company string) *lmv1.DefaultApi {
	config := lmv1.NewConfiguration()
	config.APIKey = map[string]map[string]string{
		"Authorization": map[string]string{
			"AccessID":  id,
			"AccessKey": key,
		},
	}
	config.BasePath = "https://" + company + ".logicmonitor.com/santaba/rest"

	api := lmv1.NewDefaultApi()
	api.Configuration = config

	return api
}

func main() {
  client := NewLMClient("foo", "bar", "baz")
}
```

### License
[![license](https://img.shields.io/github/license/logicmonitor/lm-sdk-go.svg?style=flat-square)](https://github.com/logicmonitor/lm-sdk-go/blob/master/LICENSE)

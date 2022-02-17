//go:build ignore
// +build ignore

package main

/*
This program generates enums.go. It can be invoked by running
go generate
*/

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

type Resource struct {
	DisplayName        string   `yaml:"displayName"`
	ShortName          []string `yaml:"shortName"`
	LongName           string   `yaml:"longName"`
	LowerCase          string   `yaml:"lowerCase"`
	ParsedResourceType []string `yaml:"parsedResourceType"`
	Apis               string   `yaml:"apis"`
	Namespaced         bool     `yaml:"namespaced"`
	PingResource       bool     `yaml:"pingResource"`
	ObjectType         string   `yaml:"objectType"`
	TitlePlural        string   `yaml:"titlePlural"`
	AdditionalHostname bool     `yaml:"additionalHostname"`
}

type Apis struct {
	PackageName string `yaml:"packageName"`
	Alias       string `yaml:"alias"`
	ApiGroup    string `yaml:"apiGroup"`
	ApiVersion  string `yaml:"apiVersion"`
}

type Resources struct {
	Resources []Resource `yaml:"resources"`
	Apis      []Apis     `yaml:"apis"`
}

const ApiRegex string = "([a-z0-9.-]+)"

// K8s core API Group
const K8sCoreAPIGroup string = "core"

var r *regexp.Regexp = regexp.MustCompile(ApiRegex)

func readConf(filename string) (*Resources, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Resources{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}

func apiExtractor(packageName string) (string, string) {
	match := r.FindAllStringSubmatch(packageName, -1)
	len := len(match)
	if len > 2 {
		return match[len-2][0], match[len-1][0]
	} else {
		panic("Invalid Match")
	}
}

func main() {
	c, err := readConf("enums.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fileName := "enums.go"

	f, err := os.Create(fileName)
	die(err)
	defer f.Close()

	for i, v := range c.Apis {
		if strings.Contains(v.PackageName, "/") {
			api, version := apiExtractor(v.PackageName)
			if v.Alias == "" {
				c.Apis[i].Alias = api + version
			}
			if v.ApiGroup == "" || api != "core" {
				c.Apis[i].ApiGroup = api
			}
			if v.ApiVersion == "" {
				c.Apis[i].ApiVersion = version
			}
		}
	}

	tmpl := template.Must(template.New("enum_template.tmpl").Funcs(lastFunc).ParseFiles("enum_template.tmpl"))

	var buf bytes.Buffer
	tmpl.Execute(&buf, struct {
		Resources Resources
		Timestamp time.Time
	}{
		Resources: *c,
		Timestamp: time.Now(),
	})

	p, err := format.Source(buf.Bytes())
	if err != nil {
		die(err)
	}

	f.Write(p)
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var lastFunc = template.FuncMap{
	"last": func(length int, i int) bool {
		return i < length-1
	},
	"notNil": func(s interface{}) bool {
		return s != nil && s != ""
	},
	"getCases": func(cases []Resource) string {
		s := ""
		for _, c := range cases {
			if len(c.Apis) == 0 {
				s = s + c.LongName + ", "
			}
		}
		return s[:len(s)-2]
	},
	"namespaced": func(cases []Resource, namespaced bool) string {
		s := ""
		for _, c := range cases {
			if namespaced == c.Namespaced {
				s = s + c.LongName + ", "
			}
		}
		return s[:len(s)-2]
	},
	"pingResource": func(cases []Resource, pingResource bool) string {
		s := ""
		for _, c := range cases {
			if pingResource == c.PingResource {
				s = s + c.LongName + ", "
			}
		}
		return s[:len(s)-2]
	},
	"additionalHostname": func(cases []Resource, additionalHostname bool) string {
		s := ""
		for _, c := range cases {
			if additionalHostname == c.AdditionalHostname {
				s = s + c.LongName + ", "
			}
		}
		return s[:len(s)-2]
	},
	"apiGroup": func(cases []Resource, apis []Apis, apiGroup string) string {
		s := ""
		m := map[string]struct{}{}
		for _, a := range apis {
			if a.ApiGroup == apiGroup {
				m[a.Alias] = struct{}{}
			}
		}
		for _, c := range cases {
			if _, ok := m[c.Apis]; ok {
				s = s + c.LongName + ", "
			}
		}
		return s[:len(s)-2]
	},
	"apiVersion": func(cases []Resource, api Apis) string {
		s := ""
		for _, c := range cases {
			if c.Apis == api.Alias {
				s = s + c.LongName + ", "
			}
		}
		return s[:len(s)-2]
	},
	"api": func(api Apis) string {
		s := api.ApiGroup
		if s == "" || s == K8sCoreAPIGroup {
			return api.ApiVersion
		} else {
			return s + "/" + api.ApiVersion
		}
	},
	"k8sObjectType": func(resource Resource) string {
		if resource.ObjectType == "" {
			return resource.DisplayName
		} else {
			return resource.ObjectType
		}
	},
	"getApiGroups": func(apis []Apis) []string {
		m := map[string]struct{}{}
		for _, v := range apis {
			m[v.ApiGroup] = struct{}{}
		}
		arr := make([]string, len(m))
		i := 0
		for key := range m {
			arr[i] = key
			i++
		}
		return arr
	},
	"apiGroupValue": func(val string) string {
		if val == "core" {
			return ""
		} else {
			return val
		}
	},
	"titlePlural": func(rt Resource) string {
		if rt.TitlePlural != "" {
			return rt.TitlePlural
		} else {
			return rt.LongName
		}
	},
}

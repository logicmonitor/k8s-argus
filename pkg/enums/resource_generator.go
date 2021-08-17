// +build ignore

package main

/*
This program generates enums.go. It can be invoked by running
go generate
*/

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
}

type Resources struct {
	Resources []Resource `yaml:"resources"`
}

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

func main() {
	c, err := readConf("enums.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fileName := "enums.go"

	f, err := os.Create(fileName)
	die(err)
	defer f.Close()

	tmpl := template.Must(template.New("enum_template.tmpl").Funcs(lastFunc).ParseFiles("enum_template.tmpl"))

	tmpl.Execute(f, struct {
		Resources Resources
		Timestamp time.Time
	}{
		Resources: *c,
		Timestamp: time.Now(),
	})
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
}

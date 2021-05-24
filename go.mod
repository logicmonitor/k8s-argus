module github.com/logicmonitor/k8s-argus

go 1.14

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/go-openapi/runtime v0.19.28
	github.com/go-openapi/strfmt v0.20.1
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/logicmonitor/k8s-collectorset-controller v2.0.0+incompatible
	github.com/logicmonitor/lm-sdk-go v2.0.0-argus3+incompatible
	github.com/prometheus/client_golang v0.9.3
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.3
	github.com/stretchr/testify v1.7.0
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
)

// https://github.com/etcd-io/etcd/issues/11749
replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

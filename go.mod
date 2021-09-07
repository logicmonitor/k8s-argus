module github.com/logicmonitor/k8s-argus

go 1.14

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/go-openapi/runtime v0.19.28
	github.com/go-openapi/strfmt v0.20.1
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.1.2
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.6 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/logicmonitor/k8s-collectorset-controller v2.0.0+incompatible
	github.com/logicmonitor/lm-sdk-go v2.0.0-argus5+incompatible
	github.com/pkg/profile v1.6.0
	github.com/prometheus/client_golang v0.9.3
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/stretchr/testify v1.7.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
)

// https://github.com/etcd-io/etcd/issues/11749
replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

NAMESPACE  := logicmonitor
REPOSITORY := collectorset-controller
VERSION    := 0.1.0-alpha.0

all: collector
	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):latest .
	docker run --rm -v $(shell pwd):/out --entrypoint=cp $(NAMESPACE)/$(REPOSITORY):latest /tmp/api.pb.go /out/api
	docker run --rm -v $(shell pwd):/out --entrypoint=cp $(NAMESPACE)/$(REPOSITORY):latest /tmp/zz_generated.deepcopy.go /out/pkg/apis/v1alpha1/
	docker tag $(NAMESPACE)/$(REPOSITORY):latest $(NAMESPACE)/$(REPOSITORY):$(VERSION)

.PHONY: collector
collector:
	cd collector && docker build -t $(NAMESPACE)/k8s-collector:latest .
	docker tag $(NAMESPACE)/k8s-collector:latest $(NAMESPACE)/k8s-collector:$(VERSION)

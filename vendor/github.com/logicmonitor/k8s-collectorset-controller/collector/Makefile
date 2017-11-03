NAMESPACE  := logicmonitor
REPOSITORY := k8s-collector
VERSION    := 0.1.0-alpha.0

all:
	docker build -t $(NAMESPACE)/$(REPOSITORY):latest .
	docker tag $(NAMESPACE)/$(REPOSITORY):latest $(NAMESPACE)/$(REPOSITORY):$(VERSION)

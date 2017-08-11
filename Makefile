NAMESPACE  := logicmonitor
REPOSITORY := argus
VERSION    := 0.1.0-alpha.3

all:
	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):latest .
	docker tag $(NAMESPACE)/$(REPOSITORY):latest $(NAMESPACE)/$(REPOSITORY):$(VERSION)

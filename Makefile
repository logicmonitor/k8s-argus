NAMESPACE  := logicmonitor
REPOSITORY := argus
VERSION    := 0.1.0-alpha.2

all:
	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):$(VERSION) .
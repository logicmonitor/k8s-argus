NAMESPACE  := logicmonitor
REPOSITORY := argus
VERSION    := 1.0.0

all:
	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):v2latest .
	docker tag $(NAMESPACE)/$(REPOSITORY):v2latest $(NAMESPACE)/$(REPOSITORY):$(VERSION)

.PHONY: docs
docs:
	cd docs/source && hugo

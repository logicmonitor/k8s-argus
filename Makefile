VERSION = 0.1.0-alpha.0

NAMESPACE := logicmonitor
NAME := argus

all: image

vendor:
	@dep ensure

.PHONY: build
build: vendor
	@mkdir -p ./build
	@env GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/logicmonitor/argus/constants.Version=$(VERSION)" -o ./build/argus

.PHONY: test
test:
	@./scripts/test.sh

.PHONY: image
image: build
	@docker build --tag $(NAMESPACE)/$(NAME):latest .

.PHONY: push
push: image
	@docker tag $(NAMESPACE)/$(NAME):latest $(NAMESPACE)/$(NAME):$(VERSION)
	@docker push $(NAMESPACE)/$(NAME):latest
	@docker push $(NAMESPACE)/$(NAME):$(VERSION)

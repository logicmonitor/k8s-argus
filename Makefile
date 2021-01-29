NAMESPACE  := logicmonitor
REPOSITORY := argus
VERSION       ?= $(shell git describe --tags --always --dirty)

default: build

lint:
ifeq ($(shell uname -s), Darwin)
	find pkg/ -type f | grep go | egrep -v "mocks|gomock" | xargs gofmt -l -d -s -w
	find pkg/ -type f | grep go | egrep -v "mocks|gomock" | xargs goimports -l -d -w
endif

build: lint

	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):$(VERSION) .

dev: lint
	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):$(VERSION) -f Dockerfile.dev .

mockgen:
	go generate ./...
test: mockgen lint
	go test ./... -v -coverprofile=coverage.txt -race
	go tool cover -html=coverage.txt -o coverage.html

devsetup:
	go get github.com/golang/mock/mockgen


.PHONY: docs
docs:
	cd docs/source && hugo

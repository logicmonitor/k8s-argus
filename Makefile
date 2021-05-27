NAMESPACE  := logicmonitor
REPOSITORY := argus
VERSION       ?= $(shell git describe --tags --always --dirty)

default: build

gofmt:
ifeq ($(shell uname -s), Darwin)
	find pkg/ -type f | grep go | egrep -v "mocks|gomock" | xargs gofmt -l -d -s -w; sync
	find pkg/ -type f | grep go | egrep -v "mocks|gomock" | xargs gofumpt -l -d -s -w; sync
	find pkg/ -type f | grep go | egrep -v "mocks|gomock" | xargs gci -w; sync
	find pkg/ -type f | grep go | egrep -v "mocks|gomock" | xargs goimports -l -d -w; sync
	find cmd/ -type f | grep go | egrep -v "mocks|gomock" | xargs gofmt -l -d -s -w; sync
	find cmd/ -type f | grep go | egrep -v "mocks|gomock" | xargs gofumpt -l -d -s -w; sync
	find cmd/ -type f | grep go | egrep -v "mocks|gomock" | xargs gci -w; sync
	find cmd/ -type f | grep go | egrep -v "mocks|gomock" | xargs goimports -l -d -w; sync
	gofmt -l -d -s -w main.go; sync
	gofumpt -l -d -s -w main.go; sync
	gci -w main.go; sync
	goimports -l -d -w main.go; sync
endif

build: gofmt
	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):$(VERSION) .

dev: gofmt
	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):$(VERSION) -f Dockerfile.dev .

lint: gofmt
	docker build --build-arg VERSION=$(VERSION) -t $(NAMESPACE)/$(REPOSITORY):$(VERSION) -f Dockerfile.lint .

mockgen:
	go generate ./...
test: mockgen lint
	go test ./... -v -coverprofile=coverage.txt -race
	go tool cover -html=coverage.txt -o coverage.html

devsetup:
	which mockgen || go get github.com/golang/mock/mockgen
	which gofumpt || go get mvdan.cc/gofumpt
	which gci || go get github.com/daixiang0/gci

.PHONY: docs
docs:
	cd docs/source && hugo

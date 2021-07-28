FROM golang:1.14 as build
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-argus
ARG VERSION
COPY ./ ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /argus -ldflags "-X \"github.com/logicmonitor/k8s-argus/pkg/constants.Version=${VERSION}\""

FROM golangci/golangci-lint:v1.40 as lint
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-argus
COPY --from=build $GOPATH/src/github.com/logicmonitor/k8s-argus ./
RUN golangci-lint run -v
# to copy file in last image, otherwise docker buildkit ignore this stage if no dependency
RUN touch /marker

FROM golang:1.14 as test
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-argus
RUN GO111MODULE=on go get github.com/golang/mock/mockgen@v1.6.0
COPY --from=build $GOPATH/src/github.com/logicmonitor/k8s-argus ./
RUN go generate ./...
# RUN chmod +x ./scripts/test.sh; sync; ./scripts/test.sh
RUN go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
RUN cp coverage.txt /coverage.txt

FROM alpine:3.6
LABEL maintainer="LogicMonitor <argus@logicmonitor.com>"
RUN apk --update add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apk/*
WORKDIR /app
COPY --from=test /coverage.txt /coverage.txt
COPY --from=lint /marker /marker
COPY --from=build /argus /bin

ENTRYPOINT ["argus"]

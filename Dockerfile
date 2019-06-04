FROM golang:1.11 as build
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-argus
ARG VERSION
COPY ./ ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /argus -ldflags "-X \"github.com/logicmonitor/k8s-argus/pkg/constants.Version=${VERSION}\""

FROM golang:1.11 as test
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-argus
RUN go get -u github.com/alecthomas/gometalinter
RUN gometalinter --install
COPY --from=build $GOPATH/src/github.com/logicmonitor/k8s-argus ./
RUN chmod +x ./scripts/test.sh; sync; ./scripts/test.sh
RUN cp coverage.txt /coverage.txt

FROM alpine:3.6
LABEL maintainer="LogicMonitor <argus@logicmonitor.com>"
RUN apk --update add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apk/*
WORKDIR /app
COPY --from=test /coverage.txt /coverage.txt
COPY --from=build /argus /bin

ENTRYPOINT ["argus"]

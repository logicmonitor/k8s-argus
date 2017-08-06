FROM golang:1.8.3 as build
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-argus
COPY ./ ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /argus

FROM golang:1.8.3 as test
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-argus
RUN go get -u github.com/alecthomas/gometalinter
RUN gometalinter --install
COPY --from=build $GOPATH/src/github.com/logicmonitor/k8s-argus ./
RUN chmod +x ./scripts/test.sh; sync; ./scripts/test.sh

FROM alpine:3.6
MAINTAINER Andrew Rynhard <andrew.rynhard@logicmonitor.com>
RUN apk --update add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apk/*
WORKDIR /app
COPY --from=build /argus /bin
ENTRYPOINT ["argus"]

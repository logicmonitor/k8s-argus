FROM alpine:3.6
MAINTAINER Andrew Rynhard <andrew.rynhard@logicmonitor.com>

RUN apk --update add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apk/*

# Watcher
WORKDIR /app
COPY ./build/argus argus

ENTRYPOINT ["./argus"]

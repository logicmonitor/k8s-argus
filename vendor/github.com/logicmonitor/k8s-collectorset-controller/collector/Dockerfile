FROM debian:buster-slim
LABEL maintainer="Andrew Rynhard <andrew.rynhard@logicmonitor.com>"
RUN apt-get -y update \
    && apt-get -y install --no-install-recommends \
    ca-certificates \
    curl \
    inetutils-traceroute \
    iputils-ping \
    ntp \
    openssl \
    perl \
    && apt-get -y clean \
    && rm -rf /var/lib/apt/lists/*

COPY entrypoint.sh /bin/
RUN chmod +x /bin/entrypoint.sh

WORKDIR /usr/local/logicmonitor
ENTRYPOINT [ "entrypoint.sh" ]

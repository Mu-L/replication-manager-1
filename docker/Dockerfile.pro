FROM golang:1.23-bookworm as builder

RUN mkdir -p /go/src/github.com/signal18/replication-manager
WORKDIR /go/src/github.com/signal18/replication-manager

COPY . .

RUN make pro cli

FROM debian:bookworm-slim

RUN mkdir -p \
        /etc/replication-manager \
        /etc/replication-manager/cluster.d \
        /var/lib/replication-manager

RUN apt-get update
RUN apt-get -y  install apt-transport-https curl
RUN mkdir -p /etc/apt/keyrings
RUN curl -o /etc/apt/keyrings/mariadb-keyring.pgp 'https://mariadb.org/mariadb_release_signing_key.pgp'

RUN echo "# MariaDB 11 Rolling RC repository list - created 2024-09-15 10:56 UTC
# https://mariadb.org/download/
X-Repolib-Name: MariaDB
Types: deb
# deb.mariadb.org is a dynamic mirror if your preferred mirror goes offline. See https://mariadb.org/mirrorbits/ for details.
# URIs: https://deb.mariadb.org/11.rc/debian
URIs: https://mirrors.ircam.fr/pub/mariadb/repo/11.rc/debian
Suites: bookworm
Components: main
Signed-By: /etc/apt/keyrings/mariadb-keyring.pgp" > /etc/apt/sources.list.d/mariadb.list

set -e;\
echo 'X-Repolib-Name: MariaDB \n\
Types: deb \n\
   URIs: https://mirrors.ircam.fr/pub/mariadb/repo/11.rc/debian \n\
   Suites: bullseye \n\
   Components: main \n\
   Signed-By: /etc/apt/keyrings/mariadb-keyring.pgp \n'\
	 > /etc/apt/sources.list.d/mariadb.list

COPY --from=builder /go/src/github.com/signal18/replication-manager/etc/local/config.toml.docker /etc/replication-manager/config.toml
COPY --from=builder /go/src/github.com/signal18/replication-manager/etc/local/masterslave/haproxy/config.toml /etc/replication-manager/cluster.d/localmasterslavehaproxy.toml
COPY --from=builder /go/src/github.com/signal18/replication-manager/etc/local/masterslave/proxysql/config.toml /etc/replication-manager/cluster.d/localmasterslaveproxysql.toml
COPY --from=builder /go/src/github.com/signal18/replication-manager/share /usr/share/replication-manager/
COPY --from=builder /go/src/github.com/signal18/replication-manager/build/binaries/replication-manager-pro /usr/bin/replication-manager
COPY --from=builder /go/src/github.com/signal18/replication-manager/build/binaries/replication-manager-cli /usr/bin/replication-manager-cli

RUN apt-get -y install mydumper ca-certificates restic mariadb-client mariadb-server mariadb-plugin-spider haproxy libmariadb-dev fuse sysbench curl
RUN curl -LO https://github.com/sysown/proxysql/releases/download/v2.5.2/proxysql_2.5.2-debian11_amd64.deb && dpkg -i proxysql_2.5.2-debian11_amd64.deb
RUN apt-get install -y adduser libfontconfig1 && curl -LO https://dl.grafana.com/oss/release/grafana_8.1.1_amd64.deb && dpkg -i grafana_8.1.1_amd64.deb
CMD ["replication-manager", "monitor", "--http-server"]
EXPOSE 10001

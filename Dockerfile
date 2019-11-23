# 1.12.6-stretch
FROM golang@sha256:35200a727dc44175d9221a6ece398eed7e4b8e17cb7f0d72b20bf2d5cf9dc00d AS build

WORKDIR /go/src/github.com/linuxuser586/cassandra

COPY . .

RUN make dist

# 3.11.5-base
FROM linuxuser586/cassandra@sha256:aee8292cb803c416e92f7d1581cc73ed0d1da270a03d83a4f8ee25297ed691ac

COPY --from=build /go/src/github.com/linuxuser586/cassandra/build/* /usr/local/bin/

COPY conf/* /etc/cassandra/

VOLUME [ "/conf" ]

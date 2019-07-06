# 1.12.6-stretch
FROM golang@sha256:35200a727dc44175d9221a6ece398eed7e4b8e17cb7f0d72b20bf2d5cf9dc00d AS build

WORKDIR /go/src/github.com/linuxuser586/cassandra

COPY . .

RUN make dist

# 3.11.4-base
FROM linuxuser586/cassandra@sha256:845de879069970ff4da96b410ce3042a0e2cf0dbd137f27c6ee359d8c51a7bfc

COPY --from=build /go/src/github.com/linuxuser586/cassandra/build/* /usr/local/bin/

COPY conf/* /etc/cassandra/

VOLUME [ "/conf" ]

# 1.12.6-stretch
FROM golang@sha256:35200a727dc44175d9221a6ece398eed7e4b8e17cb7f0d72b20bf2d5cf9dc00d AS build

WORKDIR /go/src/github.com/linuxuser586/cassandra

COPY . .

RUN make dist

# 3.11.4-base
FROM linuxuser586/cassandra@sha256:565ebe26c3d221f80188f22ce3075ef69705bbe1aec3c9a59824102e796a3dce

COPY --from=build /go/src/github.com/linuxuser586/cassandra/build/* /usr/local/bin/

COPY conf/* /etc/cassandra/

VOLUME [ "/conf" ]

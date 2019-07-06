# Distroless Cassandra (BETA)

This repository contains [Cassandra](https://cassandra.apache.org/) that runs in
a [distroless](https://github.com/GoogleContainerTools/distroless/) container for
Docker/Kubernetes

See [LICENSE](https://github.com/linuxuser586/cassandra/blob/master/LICENSE)

## Run

Only tested on Linux at this time. Use NO_TLS set to true for development. Using TLS
requires more advanced configuration and will be documented in the future. You will
need to download cqlsh, found in the Cassandra distribution, if you want to run queries.

```sh
sudo sysctl -w vm.max_map_count=1048575
docker run -it --rm -p 9042:9042 --cap-add SYS_RESOURCE --name cassandra -e NO_TLS=true linuxuser586/cassandra
```

In another terminal you can issue commands such as

```sh
docker exec -t cassandra nodetool info
```

To save data, use something like the following.

```sh
docker run -it --rm -p 9042:9042 --cap-add SYS_RESOURCE --name cassandra -v ~/cassandra/demo:/var/lib/cassandra -e NO_TLS=true linuxuser586/cassandra
```

## Configuration

You may add your own cassandra.yaml. You can start with the official version
found in the cassandra release or just add the config you are interested in
changing. For example if you just want to update the cluster you can create
a conf directory and in your cassandra.yaml in that directory add
cluster_name: 'Demo Cluster'. Then run the below command.

```sh
docker run -it --rm -p 9042:9042 --cap-add SYS_RESOURCE --name cassandra -v ~/conf:/conf -e NO_TLS=true linuxuser586/cassandra
```

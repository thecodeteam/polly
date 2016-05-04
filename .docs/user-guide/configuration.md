# Configuration

---

Polly's configuration is done with a [yaml](http://www.yaml.org/start.html) file named config.yml located in /etc/polly.

The initial install will compose a version using a mock storage driver and a simple [Bolt](https://github.com/boltdb/bolt) key value store to hold operational state.

Configuration example using Bolt DB and mock driver for non-production simple test/dev use of REST API

```
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb_test
  libstorage:
    host: tcp://localhost:7979
    profiles:
      enabled: true
      groups:
      - local=127.0.0.1
    server:
      endpoints:
        localhost:
          address: tcp://localhost:7979
      services:
        mock:
          libstorage:
            driver: mock
        mock2:
          libstorage:
            driver: mock
        vfs:
          libstorage:
            driver: vfs
```

Alternate configuration example using ScaleIO

```
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb_test
  libstorage:
    host: tcp://localhost:7981
    profiles:
      enabled: true
      groups:
      - local=127.0.0.1
    server:
      endpoints:
        localhost:
          address: tcp://localhost:7981
      services:
        mock:
          libstorage:
            driver: mock
        mock2:
          libstorage:
            driver: mock
        vfs:
          libstorage:
            driver: vfs
        scaleio:
          libstorage:
            driver: scaleio
          scaleio:
            endpoint: https://192.168.50.12/api
            insecure: true
            userName: admin
            password: Scaleio123
            systemName: cluster1
            protectionDomainName: pdomain
            storagePoolName: pool1
```
## RexRay configuration

The container host instances to be used will Polly need to have a [RexRay configuration](http://rexray.readthedocs.io/en/stable/user-guide/config/) which references the Polly server. This is an example configuration:

```
rexray:
  logLevel: debug
libstorage:
  logLevel: debug 
  host: tcp://192.0.2.123:7981
  service: scaleio
```

## Firewall port openings

The firewall on the host running Polly should be configured to open the ports for libstorage and the Polly REST API that are specified in the `\etc\polly\config.yml` file.

## Alternate key/value store option

Polly uses [libkv](https://github.com/docker/libkv) as an abstraction layer to a key value store. libkv also allows use of [Consul](https://www.consul.io/intro/getting-started/kv.html) or [Zookeeper](https://zookeeper.apache.org/doc/r3.3.3/zookeeperStarted.html) as alternatives to Bolt.

# Configuration

---

Polly's configuration is done with a [yaml](http://www.yaml.org/start.html) file named polly.yml located in /etc/polly.

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

The initial install will compose a version using a mock storage driver and a simple [Bolt](https://github.com/boltdb/bolt) key value store to hold operational state.

The firewall on the host running Polly should be configured to open the ports for libstorage and the Polly REST API that are specified in the `\etc\polly\config.yml` file.

Polly uses [libkv](https://github.com/docker/libkv) as an abstraction layer to a key value store. This allows use of [Consul](https://www.consul.io/intro/getting-started/kv.html) or [Zookeeper](https://zookeeper.apache.org/doc/r3.3.3/zookeeperStarted.html) as alternatives to Bolt.

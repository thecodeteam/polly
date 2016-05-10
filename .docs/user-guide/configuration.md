# Configuration
Setting the configuration files

---
## Main configuration file

The configuration is done with a yaml file named `config.yml`
located in /etc/polly.

An unconfigured instance will use a mock storage driver and a
simple [Bolt](https://github.com/boltdb/bolt) key value store to hold
operational state.

Configuration example using Bolt DB and mock driver for non-production simple
test/dev use of REST API

```
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb_test
  libstorage:
    host: tcp://localhost:7979
    server:
      endpoints:
        localhost:
          address: tcp://localhost:7979
      services:
        mock:
          libstorage:
            driver: mock
        vfs:
          libstorage:
            driver: vfs
```

Alternate configuration example using ScaleIO and VirtualBox.

```
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb_test
  libstorage:
    host: tcp://localhost:7981
    server:
      endpoints:
        localhost:
          address: tcp://localhost:7981
      services:
        dockerswarm_virtualbox:
          libstorage:
            driver: virtualbox
        dockerswarm_scaleio:
          libstorage:
            driver: scaleio
virtualbox:
  endpoint: http://10.0.2.2:18083
  tls: false
  volumePath: /Users/clintonkitson/VirtualBox Volumes
  controllerName: SATAController
scaleio:
  endpoint: https://192.168.50.12/api
  insecure: true
  userName: admin
  password: Scaleio123
  systemName: cluster1
  protectionDomainName: pdomain
  storagePoolName: pool1
```

## Driver configuration
Most configuration in Polly is inherited from the libStorage package. Read the
details at [libStorage](https://github.com/emccode/libStorage) to get a
better understanding of how to configure the services.

## REX-Ray configuration

`REX-Ray` is a `libStorage` compliant storage orchestrator. It can live
on the distributed hosts where the volumes must be advertised and consumed.
You can find documentation about `REX-Ray`
[here]([REX-Ray configuration](http://rexray.readthedocs.io/en/stable/user-guide/config/))

The following is an example configuration for REX-Ray talking to the Polly and
making requests for voumes on behalf of Docker. Notice how the `service`
parameters match between this and the previous configuration. This configuration
will expose `scaleio` and `virtualbox` and Volume Drivers to Docker.

```
libstorage:
  host: tcp://127.0.0.1:7981
rexray:
  modules:
    default-docker:
      type: docker
      desc: "The default docker module."
      host: "unix:///run/docker/plugins/scaleio.sock"
      libstorage:
        service: dockerswarm_scaleio
    virtualbox:
      type: docker
      desc: "The default docker module."
      host: "unix:///run/docker/plugins/virtualbox.sock"
      libstorage:
        service: dockerswarm_virtualbox
```

## Firewall port openings
The firewall on the host running Polly should be configured for the libStorage
SSL endpoint and instructions should be followed from libStorage to ensure
certificates are in order. The Polly REST API should not be exposed today.


## Alternate key/value store option

Polly uses [libkv](https://github.com/docker/libkv) as an abstraction layer to a key value store. libkv also allows use of [Consul](https://www.consul.io/intro/getting-started/kv.html) or [Zookeeper](https://zookeeper.apache.org/doc/r3.3.3/zookeeperStarted.html) as alternatives to Bolt.

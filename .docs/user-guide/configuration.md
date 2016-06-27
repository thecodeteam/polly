# Configuration
Setting the configuration files

---
## Main configuration file

The configuration is done with a yaml file named `config.yml`
located in `/etc/polly`.

Example of Polly advertising VirtualBox as `swarm_virtualbox`.

```
polly:
  host: tcp://127.0.0.1:7978
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb_test
libstorage:
  host: tcp://127.0.0.1:7981
  embedded: false
  server:
    endpoints:
      localhost:
        address: tcp://:7981
    services:
      swarm_virtualbox:
        libstorage:
          storage:
            driver: virtualbox
virtualbox:
  endpoint: http://10.0.2.2:18083
  tls: false
  volumePath: /Users/your_user/VirtualBox Volumes
  controllerName: SATAController
```

## Driver configuration
Most configuration in Polly is inherited from the libStorage package. Read the
details at [libStorage](https://github.com/emccode/libStorage) to get a
better understanding of how to configure the services.

## REX-Ray configuration

`REX-Ray` is a `libStorage` compliant storage orchestrator. It can live
on the distributed hosts where the volumes must be advertised and consumed.
You can find documentation about `REX-Ray`
[here](http://rexray.readthedocs.io/en/stable/user-guide/config/)

The following is an example configuration for REX-Ray talking to the Polly and
making requests for voumes on behalf of Docker. Notice how the `service`
parameters match between this and the previous configuration. This configuration
will expose `virtualbox` and Volume Drivers to Docker.

```
rexray:
  modules:
    default-docker:
      host:     unix:///run/docker/plugins/virtualbox.sock
      spec:     /etc/docker/plugins/virtualbox.spec
      libstorage:
        service: swarm_virtualbox
libstorage:
  host: tcp://$POLLY_IP:7981
```

## Firewall port openings
The firewall on the host running Polly should be configured for the libStorage
SSL endpoint and instructions should be followed from libStorage to ensure
certificates are in order. The Polly REST API should not be exposed today.


## Alternate key/value store option

Polly uses [libkv](https://github.com/docker/libkv) as an abstraction layer to a key value store. libkv also allows use of [Consul](https://www.consul.io/intro/getting-started/kv.html) or [Zookeeper](https://zookeeper.apache.org/doc/r3.3.3/zookeeperStarted.html) as alternatives to Bolt. Currently only Bolt, Consul, and Zookeeper are supported for back-end stores.

By default Bolt is the default backing store; however, example configurations for Consul and Zookeeper can be found below. Consul and Zookeeper instances can either be local or remote keeping in mind a network path to those services exist from Polly (i.e. firewalls, and etc). This will allow your configuration to leverage preexisting service that might exist in your application. For example, Zookeeper instances in Apache Mesos clusters.

Consul:
```
polly:
  ...
  store:
    type: consul
    endpoints: 10.50.0.1:8500
  ...
```

Zookeeper:
```
polly:
  ...
  store:
    type: zk
    endpoints: 10.50.0.1:2181
  ...
```

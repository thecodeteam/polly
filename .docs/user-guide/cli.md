# Command line interface

Polly has a single binary file `polly` that provides a CLI. This file includes
both server and client functionality. In order to leverage any functionality
outside of the `store` options, the server must be started.


## Service Control
Polly leverages the detected init software to control the service. A
`polly start` is essentially the equivalent of `service polly start`.


### Starting

Start the server in the foreground with the following command.

```shell
polly start -f
```

### Restart

Restart the server in the foreground with the following command.

```shell
polly restart
```

### Stop

Stop the server with the following command.

```shell
polly stop
```



## Volume operations

```
polly volume [command] [flags]
```

Use a `--help` at any point to receive more information.

```
Flags:
  -l, --logLevel="warn": The log level (error, warn, info, debug)
  -v, --verbose[=false]: Print verbose help information
```

### List volumes

####   Get list of volumes managed by Polly
By default this command lists only volumes that Polly has been made aware of.

```
polly volume get
```

```
$ polly volume get
- volume:
    attachments: []
    availabilityzone: zone-000
    iops: 0
    name: Volume 0
    networkname: ""
    size: 10240
    status: ""
    id: vol-000
    type: gold
    fields: {}
  volumeid: mock-vol-000
  servicename: mock
  schedulers: []
  labels:
    color: magenta
```

####   Get list of all volumes available to Polly
In order to retrieve all of the existing volumes that are available to the
Polly instance, you must specify a `--all` flag.

```
polly volume get --all
```

```
$ polly volume get --all

- volume:
    attachments: []
    availabilityzone: zone-000
    iops: 0
    name: Volume 0
    networkname: ""
    size: 10240
    status: ""
    id: vol-000
    type: gold
    fields: {}
  volumeid: mock-vol-000
  servicename: mock
  schedulers: []
  labels: {}
- volume:
    attachments: []
    availabilityzone: zone-001
    iops: 0
    name: Volume 1
    networkname: ""
    size: 40960
    status: ""
    id: vol-001
    type: gold
    fields: {}
  volumeid: mock-vol-001
  servicename: mock
  schedulers: []
  labels: {}
```

###   Offer a volume to scheduler(s)
Once have volume ID's to use, you can offer these to services or schedulers
that are attached to Polly's `libStorage` interface.

```
polly volume offer --format=json  --scheduler=<schedname1>,... \
  --volumeid=<volid>
```

```
$ polly volume offer --format=json  --scheduler=kubernetes1,mesos15 \
  --volumeid=driverName-vol-000

{"availabilityZone":"zone-000","name":"Volume 0","size":10240,"id":"vol-000",
"type":"gold","volumeid":"mock-vol-000","serviceName":"mock",
"schedulers":["kubernetes1"," mesos15"],"labels":{"color":"magenta"}}
```

###   Revoke an offer of a volume to scheduler(s)
These offers can be revoked at any time by specifying the `schedulerName`.

```
polly volume revoke --scheduler=<schedname1> --volumeid=<volid>
```


```
$ polly volume revoke --scheduler="kubernetes1" --volumeid=mock-vol-000

volume:
  attachments: []
  availabilityzone: zone-000
  iops: 0
  name: Volume 0
  networkname: ""
  size: 10240
  status: ""
  id: vol-000
  type: gold
  fields: {}
volumeid: mock-vol-000
servicename: mock
schedulers:
- 'mesos15'
labels:
  color: magenta
```

### Create label(s) on a volume
Arbitrary labels can be configured for volumes. These are set in key and value
pairs. These labels are avaialable as `fields` in the `libStorage` clients.

```
polly volume label --label=<key>=<value>,... --volumeid=<volid>
```

```
$ polly volume label --label=size=large,size2=medium --volumeid=mock-vol-000

volume:
  attachments: []
  availabilityzone: zone-000
  iops: 0
  name: Volume 0
  networkname: ""
  size: 10240
  status: ""
  id: vol-000
  type: gold
  fields: {}
volumeid: mock-vol-000
servicename: mock
schedulers:
- 'mesos15'
labels:
  color: magenta
  size: large
  size2: medium
```

### Remove label(s) on a volume
Labels can be easily removed with the following command.

```
polly volume labelremove --label=<key1>,<key2>... --volumeid=<volid>
```

```
polly volume labelremove --label=region,size --volumeid=mock-vol-000

volume:
  attachments: []
  availabilityzone: zone-000
  iops: 0
  name: Volume 0
  networkname: ""
  size: 10240
  status: ""
  id: vol-000
  type: gold
  fields: {}
volumeid: mock-vol-000
servicename: mock
schedulers:
- ' mesos15'
labels:
  color: magenta
```

###   Creates a new volume
When creating a volume you can optionally specify the scheduler and label
details.

```
polly volume create --availabilityzone=availabilityzone --iops=<int> \
 --label=<key>=<value>,... --name=name --scheduler="schedname,..." \
 --servicename=libstorage-servicename --size=<gb in int> --type=libstoragetype
```

```
$ polly volume create --servicename=mock2 --name=testing2  --size=1
volume:
  attachments: []
  availabilityzone: ""
  iops: 0
  name: testing2
  networkname: ""
  size: 1
  status: ""
  id: vol-005
  type: ""
  fields: {}
volumeid: mock2-vol-005
servicename: mock2
schedulers:
- ""
labels: {}
```

###   Removes a volume
Removing a volume is done by specifying the `volumeID`.

`polly volume --volumeid=<volid>`

```
$ polly volume remove --volumeid=mock2-vol-005
```

## Persistent Store operations
Persistent store operations provide a way to view and clear out the information
that Polly uses to track it's knowledge of volumes.


###   Get List of volumes retained in the persistent store

`polly store get`

```
$ polly store get

polly/volumeinternal/mock2-vol-006/: ""
polly/volumeinternal/mock2-vol-006/ID: mock2-vol-006
polly/volumeinternal/mock2-vol-006/Schedulers: '[""]'
polly/volumeinternal/mock2-vol-006/ServiceName: mock2
```

###   Completely erase the persistent store

***Warning: this is a destructive operation. It wipes Polly's internal
persistent store and results in non-recoverable loss of volume labels and
scheduler claims. It should not be used in a production environment.***

`polly store erase`

```
$ polly store erase

WARN[0000] erasing polly store trees                     store=&{client:0xc820148b40 boltBucket:[77 121 66 111 108 116 68 98 95 116 101 115 116] dbIndex:5 path:/tmp/boltdb timeout:10000000000 PersistConnection:false Mutex:{state:0 sema:0}}
```

## Troubleshooting
In order to troubleshoot it is suggested that you run polly in debug mode and
in the foreground. You can easily do this with `polly start -l debug -f`.

### Print Version
The version screen will display all of the commit and semantic versioning
information.

```
Binary: /go/bin/polly
SemVer: 0.1.0-dev+72+dirty
OsArch: Linux-x86_64
Branch: bugfix_store_list
Commit: c07935a9a179de29d090a03d84669408d82e0bb0
Formed: Tue, 10 May 2016 07:39:17 UTC
```

### Print the Polly environment

This can be convenient for determining the operating configuration.

`polly env`

### Determine where Polly is running

`sudo polly service status`

### Determine OS init system type

`sudo polly service initsys`

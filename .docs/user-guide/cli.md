**Command line interface**

---

**Available Commands**

### Volume operations

**Usage:** 

`polly volume [command] [flags]`

```
Flags:
  -l, --logLevel="warn": The log level (error, warn, info, debug)
  -v, --verbose[=false]: Print verbose help information

```

#### List volumes

#####   Get list of volumes managed by Polly

`polly volume get`

```
polly volume get
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

#####   Get list of all volumes known to attached storage providers, whether managed by Polly or not

`polly volume get --all`

```
polly volume get --all
 
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
- volume:
    attachments: []
    availabilityzone: zone-002
    iops: 0
    name: Volume 2
    networkname: ""
    size: 163840
    status: ""
    id: vol-002
    type: gold
    fields: {}
  volumeid: mock-vol-002
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
  volumeid: mock2-vol-001
  servicename: mock2
  schedulers: []
  labels: {}
- volume:
    attachments: []
    availabilityzone: zone-002
    iops: 0
    name: Volume 2
    networkname: ""
    size: 163840
    status: ""
    id: vol-002
    type: gold
    fields: {}
  volumeid: mock2-vol-002
  servicename: mock2
  schedulers: []
  labels: {}
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
  volumeid: mock2-vol-000
  servicename: mock2
  schedulers: []
  labels: {}
```

####   Offer a volume to scheduler(s)

`polly volume offer --format=json  --scheduler="<schedname1>,..." --volumeid="<volid>"`
 
 ```
polly volume offer --format=json  --scheduler="kubernetes1, mesos15" --volumeid=mock-vol-000

{"availabilityZone":"zone-000","name":"Volume 0","size":10240,"id":"vol-000","type":"gold","volumeid":"mock-vol-000","serviceName":"mock","schedulers":["kubernetes1"," mesos15"],"labels":{"color":"magenta"}}
 ```

####   Revoke an offer of a volume to scheduler(s)

`polly volume revoke --scheduler=<schedname1> --volumeid=<volid>`

```
polly volume revoke --scheduler="kubernetes1" --volumeid=mock-vol-000
 
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
 
#### Create label(s) on a volume

`polly volume label --label="<key>=<value>,..." --volumeid=<volid>`

```
polly volume label --label="size=large" --volumeid=mock-vol-000

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
  size: large
```

#### Remove label(s) on a volume

`polly volume labelremove --label="<key1>,<key2>..." --volumeid=<volid>`

```
polly volume labelremove --label="region,size" --volumeid=mock-vol-000
 
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

####   Creates a new volume

`polly volume create --availabilityzone="availabilityzone" --iops=<int> --label="<key>=<value>,..." --name="name" --scheduler="schedname,..." --servicename="libstorage-servicename" --size=<int> --type="libstoragetype"`

```
polly volume create --servicename=mock2 --name=testing2  --size=1
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

####   Removes a volume
 
`polly volume --volumeid=<volid>`

```
polly volume remove --volumeid=mock2-vol-005
``` 
 
### Persistent Store operations

**Usage:**

`polly store [command] [flags]`

####   Get List of volumes retained in the persistent store

`polly store get`

```
polly store get

polly/volumeinternal/mock2-vol-006/: ""
polly/volumeinternal/mock2-vol-006/ID: mock2-vol-006
polly/volumeinternal/mock2-vol-006/Schedulers: '[""]'
polly/volumeinternal/mock2-vol-006/ServiceName: mock2
```

####   Completely erase the persistent store

***Warning: this is a destructive operation. It wipes Polly's internal persistent store and results in non-recoverable loss of volume labels and scheduler claims. It should not be used in a production environment.***

`polly store erase`

```
polly store erase
 
WARN[0000] erasing polly store trees                     store=&{client:0xc820148b40 boltBucket:[77 121 66 111 108 116 68 98 95 116 101 115 116] dbIndex:5 path:/tmp/boltdb timeout:10000000000 PersistConnection:false Mutex:{state:0 sema:0}}
```

### Troubleshooting

#### Print Version

`polly version`


#### Print the Polly environment

This can be convenient for determining the operating configuration

`polly env`

#### Determine where Polly is running

`sudo polly service status`

#### Determine OS init system type

`sudo polly service initsys`
.

#!/bin/bash
echo "populate volumes"
 curl -d '{"service": "mock", "name": "mysql-vol", "volumeType": "nas", "size": 10240, "iops": 135, "availabilityZone": "east", "schedulers": ["mesos-15"]}' -i 127.0.0.1:7980/admin/volumes

 curl -d '{"service": "mock", "name": "postgreSQL-vol", "volumeType": "nas", "size": 10240, "iops": 1200, "availabilityZone": "west", "schedulers": ["kubernetes-1"]}' -i 127.0.0.1:7980/admin/volumes

echo -e ""
echo -e "listing all volumes"
curl -i localhost:7980/admin/volumes
echo -e ""

echo -e "delete volume"
curl -X DELETE -i localhost:7980/admin/volumes/mock-vol-004

echo -e "get all volumes after deleting one"
curl -i localhost:7980/admin/volumes

echo -e "associate a volume with a scheduler"
curl -d '{"volumeID": "mock-vol-000", "schedulers": [ "mesos-99"]}' -i 127.0.0.1:7980/admin/volumeoffer

echo -e ""

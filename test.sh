#!/bin/bash
echo "populate volumes"
curl -d '{"Scheduler":"Kubernetes-1", "ServiceName": "mock", "Name": "myAWSvol5"}' -i localhost:8080/admin/volumes
curl -d '{"Scheduler":"Swarm-1", "ServiceName": "mock", "Name": "myAWSvol2016"}' -i localhost:8080/admin/volumes
echo -e ""
echo -e "listing volumes"
curl -i localhost:8080/admin/volumes
echo -e ""

echo -e "delete volume"
curl -X DELETE -i localhost:8080/admin/volumes/mock-vol-005

echo -e "get all volumes"
curl -i localhost:8080/admin/volumes

echo -e ""

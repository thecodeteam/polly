# polly - Polymorphic Volume Scheduling

![polly](docs/images/Polly the Parrot_Containers.png)

## Storage scheduling for container runtimes

Polly implements a centralized storage scheduling service that will work with popular container schedulers of different application platforms for containerized workloads. Polly is an open source platform that supports use of external storage, with scheduled containerized workloads, at scale.

## How it works

Polly creates an abstraction layer between the container runtime and the storage infrastructure. I can be used for creating, mapping, snapshotting and deleting persistent data volumes on a multitude of storage platforms.

## Platforms
Planned supported container runtime platforms:
 - Docker Swarm
 - Mesos
 - Kubernetes
 - Cloud Foundry

Storage platforms:
 - AWS EC2 (EBS)
 - Google Compute Engine
 - EMC ScaleIO
 - EMC XtremIO
 - EMC VMAX
 - EMC Isilon
 - OpenStack Cinder
 - VirtualBox

# polly - Polymorphic Volume Scheduling

![polly](docs/images/Polly the Parrot_Containers.png)

## Storage scheduling for container runtimes

Polly implements a centralized storage scheduling service that will work with popular container schedulers of different application platforms for containerized workloads. Polly is an open source platform that supports use of external storage, with scheduled containerized workloads, at scale. It can be used for creating, mapping, snapshotting and deleting persistent data volumes on a multitude of storage platforms.

## What it does

Container runtime schedulers need to be integrated with every aspect of available hardware resources, including persistent storage. When requesting resources for an application the scheduler gets offers for CPU, RAM _and_ disk.

To be able to offer persistent storage for applications the scheduler needs to know about the underlying storage infrastructure. Having separate modules for every type of storage infrastructure is inefficient, and that is why Polly exists. It creates an abstraction layer to support multiple storage infrastructure layers for multiple containers runtime schedulers.

## Example workflow

1. An application requires highly available storage with a specific set of policies applied
1. The scheduler receives a request to start the application
3. The scheduler checks with Polly to see if there is any available storage offers that matches the requirements
4. Polly requests the volume(s) to be mapped to the container
5. Scheduler starts the container with persistent storage

## Platforms
Planned supported container runtime schedulers:
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

## GUI

1. Clone this repo
2. from the `gui` folder type: `npm install && npm start`
3. Go to `http://localhost:8000/volumes.html`

There may be problems with CORS on the server-side function. For now, enable COORS on your browser. [Chrome CORS Extension](https://chrome.google.com/webstore/detail/allow-control-allow-origi/nlfbmbojpeacfghkpbjhddihlkkiljbi?hl=en)

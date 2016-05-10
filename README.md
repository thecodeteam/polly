# polly - Polymorphic Volume Scheduling
[![Build Status](https://travis-ci.org/emccode/polly.svg?branch=master)](https://travis-ci.org/emccode/polly)
[![Go Report Card](http://goreportcard.com/badge/emccode/polly)](http://goreportcard.com/report/emccode/polly) [![codecov.io](https://codecov.io/github/emccode/polly/coverage.svg?branch=master)](https://codecov.io/github/emccode/polly?branch=master) [![Download](http://api.bintray.com/packages/emccode/polly/stable/images/download.svg)](https://dl.bintray.com/emccode/polly/)
[![Download](http://api.bintray.com/packages/emccode/polly/staged/images/download.svg)](https://dl.bintray.com/emccode/polly/) [![Docs](https://readthedocs.org/projects/polly-scheduler/badge/?version=latest)](http://polly-scheduler.readthedocs.io/en/latest/?badge=latest)

![polly](https://raw.githubusercontent.com/emccode/polly/master/.docs/images/polly.png)

## Storage scheduling for container schedulers
`Polly` implements a centralized storage scheduling service that integrates with popular `container schedulers` of different application platforms for containerized workloads. It is an open source framework that supports use of external storage, with scheduled containerized workloads, at scale. It can be used to centralize the control of creating, mapping, snapshotting and deleting persistent data volumes on a multitude of storage platforms.

## Full Docuemntation
Continue reading the full documentation at [ReadTheDocs](http://polly-scheduler.readthedocs.io/en/latest/).

## Key Features
- Centralized control and distribution of storage resources
- Offer based mechanism for advertising storage to container schedulers
- Framework supporting direct integration to any container scheduler, storage orchestrator, and storage platform

## What it does
Container runtime schedulers need to be integrated with every aspect of available hardware resources, including persistent storage. When requesting resources for an application the scheduler gets offers for CPU, RAM _and_ disk.

To be able to offer persistent storage in a scalable way, the application and container scheduler needs awareness of the available resources from the underlying storage infrastructure.

## Example workflow

1. An application requires highly available storage with a specific set of
policies applied.
1. The scheduler receives a request to start the application.
3. The scheduler checks with Polly or already is aware of outstanding offers
for storage resources.
4. Scheduler send request to container runtimes to start the container with
persistent storage.
5. Container runtime requests volume access from Polly.
6. Container runtime orchestrates process of starting container and attaching
persistent storage with help from a libStorage storage orchestrator.

## Framework to support the following platforms
Polly provides an open framework to enable integration to any container, cloud, or storage platform.

Container runtime schedulers:
 - Docker Swarm
 - Mesos
 - Kubernetes
 - Cloud Foundry

Cloud platforms:
- AWS EC2 (EBS)
- Google Compute Engine
- OpenStack
 - Private Cloud
 - Public Cloud (RackSpace, and others)

Planned supported storage platforms:
 - EMC ScaleIO
  - XtremIO
  - VMAX
  - Isilon
 - Others
 - VirtualBox

## Installation
The following command will install Polly.  If using
`CentOS`, `RedHat`, `Ubuntu`, or `Debian` the necessary service manager is used
to bootstrap the process on startup.  

`curl -sSL https://dl.bintray.com/emccode/polly/install | sh -s stable`

You can also install the latest staged release with the following command.

`curl -sSL https://dl.bintray.com/emccode/polly/install | sh -s staged`


## libStorage
Polly makes use of the open source storage plugin framework [libStorage](https://github.com/emccode/libstorage) to enable storage orchestrator tools and container runtimes to make requests of storage. Any storage platform that has a driver implementation for the libStorage framework will work with Polly.

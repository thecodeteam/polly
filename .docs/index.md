# polly
Polymorphic Volume Scheduling

![polly](images/polly.png)

## Overview
`Polly` implements a centralized storage scheduling service that integrates with popular `container schedulers` of different application platforms for containerized workloads. It is an open source framework that supports use of external storage, with scheduled containerized workloads, at scale. It can be used to centralize the control of creating, mapping, snapshotting and deleting persistent data volumes on a multitude of storage platforms.

## Key Features
- Configure services to interface with schedulers and groups of container runtimes
- Restrict existing volumes and new volumes to services
- Centralized control and distribution of storage resources
- Offer based mechanism for advertising storage to container schedulers
- Framework supporting direct integration to any container scheduler, storage orchestrator, and storage platform
- Basic Volume Management capabilities

## What it does
Container runtime schedulers need to be integrated with every aspect of available hardware resources, including persistent storage. When requesting resources for an application the scheduler gets offers for CPU, RAM _and_ disk.

To be able to offer persistent storage in a scalable way, the application and container scheduler needs awareness of the available resources from the underlying storage infrastructure.

## Example workflow

1. An application requires highly available storage with a specific set of policies applied
1. The scheduler receives a request to start the application
3. The scheduler checks with Polly or already has an off from Polly for storage resources
4. Polly requests the volume(s) to be mapped to the container
5. Scheduler issues request to start the container with persistent storage
6. Container runtime orchestrates process of starting container and attaching persistent storage

## Container runtime scheduler support
 - Docker Swarm
 - Mesos
 - Kubernetes
 - Cloud Foundry

## Cloud platform support
- AWS EC2 (EBS)
- Google Compute Engine
- OpenStack
 - Private Cloud
 - Public Cloud (RackSpace, and others)

## Storage platform support
 - EMC ScaleIO
 - XtremIO
 - VMAX
 - Isilon
 - VirtualBox
 - Others

## libStorage
Polly makes use of the open source storage plugin framework [libStorage](https://github.com/emccode/libstorage) to enable storage orchestrator tools and container runtimes to make requests of storage. Any storage platform that has a driver implementation for the libStorage framework will work with Polly.

### Hello Polly
In the grand tradition of technical documentation, the first true end-to-end
example of Polly is called `Hello Polly`. It showcases a two-node
deployment with the first node configured with REX-Ray talking to a
Polly/libstorage server and the second node as merely a REX-Ray client
talking to the Polly server on the first node. Both nodes have Docker (1.11+)
installed and configured to leverage Polly for persistent storage.

The below example does have a few requirements:

 * VirtualBox 5.0+
 * Vagrant 1.8+
 * Ruby 2.0+

#### Start Polly Vagrant Environment
Before bringing the Vagrant environment online, please ensure it is
accomplished in a clean directory:

```sh
$ cd $(mktemp -d)
```

Inside the newly created, temporary directory, download the Polly
[Vagrantfile](https://github.com/emccode/polly/blob/master/Vagrantfile):

```sh
$ curl -fsSLO https://raw.githubusercontent.com/emccode/polly/master/Vagrantfile
```

Now it is time to bring the Polly environment online:

!!! note "note"

    The next step could potentially open up the system on which the command
    is executed to security vulnerabilities. The Vagrantfile brings the
    VirtualBox web service online if it is not already running. However,
    in the name of simplicity the Vagrantfile also disables the web server's
    authentication module. Please do not disable authentication for the
    VirtualBox web server if this example is being executed on an open network
    or without some type of firewall in place.

```sh
$ vagrant up
```

Once the command has been completed successfully there will be two VMs online
named `node0` and `node1`. Both nodes are running Docker and REX-Ray; however,
`node0` has Polly configured to act as a libStorage server. Both REX-Ray instances
will be talking to Polly for managing volumes for your container
runtimes (i.e. Docker).

Now that the environment is online it is time to showcase Docker leveraging
REX-Ray to create persistent storage as well as illustrating REX-Ray's
distributed deployment capabilities.

#### Node 0
First, SSH into `node0`

```sh
$ vagrant ssh node0
```

From `node0` use Docker with REX-Ray backed by a Polly server to create
a new volume named `hellopersistence`:

```sh
vagrant@node0:~$ docker volume create --driver rexray --opt size=1 \
                 --name hellopersistence
```

You can verify that REX-Ray provisioned the volume by running the following
command:

```sh
vagrant@node0:~$ rexray volume
```

Since the volume creation was actually created via a Polly server, you can
verify that Polly is tracking the volumes created via the REX-Ray interface:

```sh
vagrant@node0:~$ polly volume ls
```

After the volume is created, mount it to the host and container using the
`--volume-driver` and `-v` flag in the `docker run` command:

```sh
vagrant@node0:~$ docker run -tid --volume-driver=rexray \
                 -v hellopersistence:/mystore \
                 --name temp01 busybox
```

Create a new file named `myfile` on the file system backed by the persistent
volume using `docker exec`:

```sh
vagrant@node0:~$ docker exec temp01 touch /mystore/myfile
```

Verify the file was successfully created by listing the contents of the
persistent volume:

```sh
vagrant@node0:~$ docker exec temp01 ls /mystore
```

Remove the container that was used to write the data to the persistent volume:

```sh
vagrant@node0:~$ docker rm -f temp01
```

Finally, exit the SSH session to `node0`:

```sh
vagrant@node0:~$ exit
```

#### Node 1
It's time to connect to `node1` and use the volume `hellopersistence` that was
created in the previous section from `node0`.

!!! note "note"

    While `node1` runs both the Docker and REX-Ray services like `node0`, the
    REX-Ray service on `node1` in no way understands or is configured for the
    VirtualBox storage driver. All interactions with the VirtualBox web service
    occurs via `node0`'s Polly server with which `node1` communicates.

Use the vagrant command to SSH into `node1`:

```sh
$ vagrant ssh node1
```

Next, create a new container that mounts the existing volume,
`hellopersistence`:

```sh
vagrant@node1:~$ docker run -tid --volume-driver=rexray \
                 -v hellopersistence:/mystore \
                 --name temp01 busybox
```

The next command validates the file `myfile` created from `node0` in the
previous section has persisted inside the volume across machines:

```sh
vagrant@node1:~$ docker exec temp01 ls /mystore
```

Finally, exit the SSH session to `node1`:

```sh
vagrant@node1:~$ exit
```

#### Cleaning Up
Be sure to kill the VirtualBox web server with a quick `killall vboxwebsrv` and
to tear down the Vagrant environment with `vagrant destroy -f`. Omitting these
commands will leave the web service and REX-Ray/Polly Vagrant nodes online and
consume additional system resources.

#### Congratulations
REX-Ray with the use of Polly on the backend has been used to provide
persistence for stateless containers!

## Getting Help
Having issues? No worries, let's figure it out together.

### GitHub and Slack
If a little extra help is needed, please don't hesitate to use
[GitHub issues](https://github.com/emccode/polly/issues) or join the active
conversation on the
[EMC {code} Community Slack Team](http://community.emccode.com/) in
the #project-polly channel

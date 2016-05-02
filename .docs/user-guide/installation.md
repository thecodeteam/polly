# Installation

## Overview
Polly is written in Go, so there are typically no dependencies that must be installed alongside its single binary file. The manual methods can be extremely simple through tools like `curl`. You also have the opportunity to perform install steps individually. Following the manual installs, [configuration](./config.md) must take place.

---

### Simple default install via curl

The following command will install the most recent, stable build of Polly to:
` /usr/bin/polly`.

```shell
sudo curl -sSL https://dl.bintray.com/emccode/polly/install | sh -s stable
```

### Installing a specific pre-built binary version
There are a handful of necessary manual steps to properly install Polly
from pre-built binaries.

1. Download the proper binary. There are also pre-built binaries available for
the various release types.

    Version  | Description
    ---------|------------
    [Unstable](https://dl.bintray.com/emccode/polly/unstable/latest/) | The most up-to-date, bleeding-edge, and often unstable Polly binaries.
    [Stable](https://dl.bintray.com/emccode/polly/stable/latest/) | The most up-to-date, stable Polly binaries.

2. Uncompress and move the binary to the proper location. `/usr/bin`
is the preferred normal location, but this path is not required.

### Build and install from source

`Polly` is also fairly simple to build from source:

#### Build using a Docker container

```shell
SRC=$(mktemp -d 2> /dev/null || mktemp -d -t polly 2> /dev/null) && cd $SRC && docker run --rm -it -e GO15VENDOREXPERIMENT=1 -v $SRC:/usr/src/polly -w /usr/src/polly golang:1.5.1 bash -c "git clone https://github.com/emccode/polly.git -b master . && make build-all”
```

#### Conventional build from source

If you'd prefer to not use `Docker` to build `Polly,` [Install go 1.6](https://golang.org/doc/install) or later and set the GOPATH environment variable. 

Also set the environment variable GO15VENDOREXPERIMENT:

```shell
export GO15VENDOREXPERIMENT=1
```

##### clone the polly repo

```shell
mkdir -p ~/work/src/github.com/emccode
cd ~/work/src/github.com/emccode
git clone https://github.com/emccode/polly.git

# change directories into the freshly-cloned repo
cd polly

# build polly
make build-all
```

After either of the above methods for building `polly` there should be a `.bin` directory in the current directory, and inside `.bin` will be binaries for Linux-i386, Linux-x86-64, and Darwin-x86-64.

```shell
$ ls .bin/*/polly
-rwxr-xr-x. 1 root 14M Sep 17 10:36 .bin/Darwin-x86_64/polly*
-rwxr-xr-x. 1 root 12M Sep 17 10:36 .bin/Linux-i386/polly*
-rwxr-xr-x. 1 root 14M Sep 17 10:36 .bin/Linux-x86_64/polly*
```
copy the binary to `/usr/bin`

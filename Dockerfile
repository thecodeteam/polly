FROM ubuntu:14.04
MAINTAINER <EMC{code}>

# To build this dockerfile first ensure that it is named "Dockerfile"
# make sure that a directory "docker_resources" is also present in the same directory as "Dockerfile",
#   and that "docker_resources" contains the files "go-wrapper" and "get_go-bindata_md5.sh"

# Assuming:
# Your dockerhub username: dhjohndoe
# Your github username: ghjohndoe
# Your polly fork is checked out in $HOME/go/src/github.com/ghjohndoe/polly/

# To build a Docker image using this Dockerfile:
# docker build -t dhjohndoe/golang-glide:0.1.0 .

# To build polly using this Docker image:
# docker pull dhjohndoe/golang-glide:0.1.0
# If cutting and pasting the next line remove '\' and '#' characters remember to replace ghjohndoe and dhjohndoe
# docker run -v $HOME/go/src/github.com/ghjohndoe/polly/:/go/src/github.com/emccode/polly/ \
# -v $HOME/build/polly/pkg/:/go/pkg/ \
# -v $HOME/build/polly/bin/:/go/bin/ \
# -w=/go/src/github.com/emccode/polly/ dhjohndoe/golang-glide:0.1.0

# after building build resources will be placed in $HOME/build/polly/ in the pkg/ and bin/ directories
# to run an instance of polly inside of a Docker container perform another docker run, this time entering the shell
# If cutting and pasting the next line remove '\' and '#' characters and remember to replace ghjohndoe and dhjohndoe
# docker run -ti -v $HOME/go/src/github.com/ghjohndoe/polly/:/go/src/github.com/emccode/polly/ \
# -v $HOME/build/polly/pkg/:/go/pkg/ \
# -v $HOME/build/polly/bin/:/go/bin/ \
# -w=/go/src/github.com/emccode/polly/ dhjohndoe/golang-glide:0.1.0 /bin/bash

# once inside run-
#   polly service start -f
# note that no configuration file exists in /etc/polly
# also note that you must terminate the instance manually with ctrl+c and exit

RUN apt-get update && apt-get install -y --no-install-recommends software-properties-common
RUN add-apt-repository ppa:masterminds/glide

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
        git \
        curl \
        g++ \
        gcc \
        libc6-dev \
        make \
        glide \
    && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.6.2
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 e40c36ae71756198478624ed1bb4ce17597b3c19d243f3f0899bb5740d56212a

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

ADD docker_resources/go-wrapper /usr/local/bin/
ADD docker_resources/get_go-bindata_md5.sh /home/

EXPOSE 7978
EXPOSE 7979

# The get_go-bindata_md5.sh script is required to resolve build errors related to go-bindata and akutz's md5checksum
CMD ["/bin/bash", "-c", "/home/./get_go-bindata_md5.sh; make clean; make deps; make"]

#! /usr/bin/env bash

make deps GOARCH=amd64 GOOS=linux
make GOARCH=amd64 GOOS=linux

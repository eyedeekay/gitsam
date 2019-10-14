#! /usr/bin/env sh

GO111MODULE=off
GOPATH=$(pwd)/go

cd go/src/github.com/eyedeekay
mkdir -p ../.gitsam_secure/
$GOPATH/bin/gitsam 2>&1 | tee gitsam.log

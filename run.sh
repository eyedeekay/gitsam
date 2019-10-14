#! /usr/bin/env sh

GO111MODULE=off
GOPATH=$(pwd)/go

cd go/src/github.com/eyedeekay
go get -u github.com/eyedeekay/gitsam/gitsam
cp ~/.ssh/github_id_rsa.pub ./id_rsa.pub
mkdir -p ../.gitsam_secure/
$GOPATH/bin/gitsam 2>&1 | tee gitsam.log

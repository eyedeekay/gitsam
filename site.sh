#! /usr/bin/env sh

GH="github.com/"
UN="eyedeekay"
GO111MODULE=off
GOPATH=$(pwd)/go
mkdir -p $GOPATH
GIT_REPOS="$GH$UN/gitsam/gitsam
$GH$UN/eephttpd/eephttpd
$GH$UN/sam-forwarder/samcatd
$GH$UN/goSam
$GH$UN/httptunnel/httpproxy
$GH$UN/httptunnel/multiproxy/browserproxy
$GH$UN/udptunnel
$GH$UN/checki2cp
$GH$UN/sam3
$GH$UN/outproxy/outproxy"

for x in $GIT_REPOS; do
    go get -u $x
done

./run.sh
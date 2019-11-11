#! /usr/bin/env sh

export SCRIPTDIR=$(dirname $(readlink -f "$0"))

. $SCRIPTDIR/.config
export CONTENT=$SCRIPTDIR/$UN

GIT_REPOS="$GH$UN/gitsam
$GH$UN/eephttpd
$GH$UN/sam-forwarder
$GH$UN/goSam
$GH$UN/httptunnel
$GH$UN/udptunnel
$GH$UN/checki2cp
$GH$UN/sam3
$GH$UN/outproxy
$GH$UN/i2pdig
$GH$UN/i2psetproxy.js
$GH$UN/i2pbutton
$GH$UN/accessregister
$GH$UN/defcon
$GH$UN/i2p-tools-1
$GH$UN/go-i2cp
$GH$UN/firefox.profile.i2p
$GH$UN/apt-transport-i2phttp
$GH$UN/apt-transport-i2p
$GH$UN/ramp
$GH$UN/Jsam
$GH$UN/basic-tunnel-tutorial
$GH$UN/geti2p64
$GH$UN/So-You-Want-To-Write-A-SAM-Library
$GH$UN/i2pkeys
$GH$UN/i2p-ssh-config"

mkdir -p "$CONTENT" $SCRIPTDIR/.gitsam_secure
cd "$CONTENT" && pwd

for x in $GIT_REPOS; do
    THEDIR=$CONTENT/$(echo "$x" | sed "s|$GH$UN/||g")
    echo $THEDIR
    mkdir -p $THEDIR
    git clone --mirror "$x" 2>&1 | grep -v fatal
    cd $THEDIR
    pwd && ls "$THEDIR.git"
    git init --separate-git-dir=$THEDIR.git
    git checkout -f
    git checkout -b pages
    git fetch --force --all --keep --progress --update-head-ok --tags
    git update-server-info -f && echo "updated server info"
    git fsck
    cd "$CONTENT"
done

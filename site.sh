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
$GH$UN/outproxy"

mkdir -p "$CONTENT" $SCRIPTDIR/.gitsam_secure
cd "$CONTENT" && pwd

for x in $GIT_REPOS; do
    THEDIR=$CONTENT/$(echo "$x" | sed "s|$GH$UN/||g")
    echo $THEDIR
    mkdir -p $THEDIR
    git clone --mirror "$x" 2>&1 | grep -v fatal
    cd $THEDIR
    pwd && echo "$CONTENT/$THEDIR.git"
    git init --separate-git-dir="$THEDIR.git"
    git checkout -f
    git fetch --force --all --keep --progress --update-head-ok --tags
    git update-server-info -f && echo "updated server info"
    git fsck
    cd "$CONTENT"
done

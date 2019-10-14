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
cd "$CONTENT"

for x in $GIT_REPOS; do
    THEDIR=$CONTENT/$(echo "$x" | sed "s|$GH$UN/||g")
    mkdir -p $(echo "$x" | sed "s|$GH$UN/||g")
    git clone "$x" $THEDIR/ 2>&1 | grep -v fatal
    cd $THEDIR
    pwd
    git fetch --force --all --keep --update-head-ok --progress
    git update-server-info -f && echo "updated server info"
    cd "$CONTENT"
done

$SCRIPTDIR/run.sh

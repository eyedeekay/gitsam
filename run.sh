#! /usr/bin/env sh

export SCRIPTDIR=$(dirname $(readlink -f "$0"))

. $SCRIPTDIR/.config
export CONTENT=$SCRIPTDIR/$UN

cd "$CONTENT" && $SCRIPTDIR/gitsam/gitsam
#! /usr/bin/env sh

export SCRIPTDIR=$(dirname $(readlink -f "$0"))

. $SCRIPTDIR/.config
export CONTENT=$SCRIPTDIR/$UN

cd "$CONTENT" && $SCRIPTDIR/gitsam/gitsam -il 1 -ol 1 -iq 8 -oq 8 -ib 3 -ob 3 -pk ../.gitsam_secure/id_rsa.pub
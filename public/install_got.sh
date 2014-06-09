#!/bin/bash -
#
# install got
# 2014-06-09 by codeskyblue
#

# example: check got go get github.com/gobuild/got
fix(){
	if ! which "$1" &>/dev/null
	then
		echo "fix $1"
		shift
		echo "$@"
		"$@"
	fi
	return 0
}

die(){
	echo "[DIE] $@"
	exit 1
}

test -n "$GOPATH" || die "need go installed"
fix got go get github.com/gobuild/got || die "got install failed"

if test $# -ne 0
then
	got "$@"
fi

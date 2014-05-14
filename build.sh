#!/bin/bash -
#
#
set -eu
(cd pkg/config && go install)
(cd pkg/xrpc && go install)
(cd routers && go install)
(cd models && go install)
go build

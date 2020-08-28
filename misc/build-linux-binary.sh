#!/bin/bash

# update and install require packages
apt-get -qq update
apt-get -qqy install gox

# set environments
XC_ARCH=amd64
XC_NAME=ltsv-analyzer
VERSION=${1:-$(cat version)}
REVISION=${2:-$(git describe --always)}
GOVERSION=$(go version)
BUILDDATE=$(date '+%Y/%m/%d %H:%M:%S %Z')
ME=$(whoami)

echo "build linux binary $XC_NAME"
GO111MODULE=on gox -os "linux" -arch "$XC_ARCH" -osarch "!darwin/arm" -ldflags "-X main.version=$VERSION -X main.revision=$REVISION -X \"main.goversion=$GOVERSION\" -X \"main.builddate=$BUILDDATE\" -X \"main.builduser=$ME\"" -output "pkg/{{.OS}}_{{.Arch}}/$XC_NAME" .

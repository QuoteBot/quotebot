#! /bin/bash

outputdir=bin/

function buildRelease() {
    echo build linux amd64
    GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bin/ .
}

function buildForDebug {
    echo "build for debug on host"
    go build -race -o bin/ .
}

[ -d $outputdir ] || mkdir $outputdir
[ "$1" == "debug" ] && buildForDebug || [ "$1" == "release" ] && buildRelease || buildForDebug
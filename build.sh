#!/usr/bin/env bash

GOPATH=$(pwd)

export GOPATH=$GOPATH

go build -o bin/hrentabd src/*.go;

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            echo "HELP!"
            exit 0
            ;;
        -r|--run)
            bin/hrentabd
            ;;
    esac
done
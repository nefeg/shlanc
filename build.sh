#!/usr/bin/env bash

GOPATH=$(pwd)

export GOPATH=$GOPATH
export GOBIN=$GOPATH/bin

go install hrentabd ; # [SH]lanc [L]ike [A]s [N]ot [C]ron

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            echo "HELP!"
            exit 0
            ;;

        -r|--run)
            `$GOBIN/hrentabd $GOPATH/config.json`
            ;;
    esac
done
echo "Done"
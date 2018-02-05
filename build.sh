#!/usr/bin/env bash

GOPATH=$(pwd)

export GOPATH=$GOPATH
export GOBIN=$GOPATH/bin

go install shlancd shlanc; # [SH]lanc [L]ike [A]s [N]ot [C]ron

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            echo "HELP!"
            exit 0
            ;;

        -r|--run)
            `$GOBIN/shlancd $GOPATH/config.json`
            ;;
    esac
done
echo "Done"
#!/bin/bash

set -e
dir_name=$(dirname $(readlink -m $0))
pushd $dir_name
    cat ../upload/$1
popd

# vim:ai:et:sts=4:tw=80:sw=4:

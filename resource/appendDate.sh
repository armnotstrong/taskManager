#!/bin/bash
set -e
dir_name=$(dirname $(readlink -m $0))
pushd $dir_name
    cp ../upload/$1 ../download/$1
    echo $(date) >> ../download/$1
    mv ../download/$1 ../download/$2
popd


#!/bin/bash
script_dir="$(dirname "$0")"

pushd $script_dir/src/server
# go run .
weaver multi deploy weaver.toml
popd

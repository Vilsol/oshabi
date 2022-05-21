#!/usr/bin/env bash

set -ex

pushd ../../
pwd
go run ./dll-copy "$1"
popd

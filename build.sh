#!/bin/sh

ZN=zinc
ZN_DEV=zinc-devtool
ZN_SERVER=zinc-server
ZN_PLAYGROUND=zinc-playground

# build Zn
rm -f ./$ZN ./$ZN_DEV ./$ZN_SERVER ./$ZN_PLAYGROUND

echo '=== build [Zn] ==='
go build -o zinc ./cmd/$ZN

echo '=== build [Zn-devtool] ==='
go build -o zinc-devtool ./cmd/$ZN_DEV

echo '=== build [Zn-server] ==='
go build -o zinc-server ./cmd/$ZN_SERVER

echo '=== build [Zn-playground] ==='
go build -o zinc-playground ./cmd/$ZN_PLAYGROUND

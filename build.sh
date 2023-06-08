#!/bin/sh

ZN=zinc
ZN_DEV=znt

# build Zn
rm -f ./$ZN ./$ZN_DEV

echo '=== build [Zn] ==='
go build -o zinc

echo '=== build [Zn-devtool] ==='
go build -o zinc-devtool ./cmd/zinc-devtool

echo '=== build [Zn-server] ==='
go build -o zinc-server ./cmd/zinc-server

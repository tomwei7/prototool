#!/bin/bash
set -ex
GOOS=windows go build -o prototool-windows-amd64.exe ./cmd/prototool
GOOS=linux  go build -o prototool-linux-amd64 ./cmd/prototool
GOOS=darwin go build -o prototool-darwin-amd64 cmd/prototool
for f in $(ls | grep prototool-)
do
    tar -Jcf ${f}.tar.xz ${f}
done
rm prototool-linux-amd64 prototool-darwin-amd64 prototool-windows-amd64.exe

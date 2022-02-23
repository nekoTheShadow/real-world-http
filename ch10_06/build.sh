#!/bin/bash

# WSL2で開発しているが、open.StartがWSL2ではうまくうごかないため、
# Windows用のバイナリを作って、cmdで動かす

set -x
rm -v ch10_06.exe
GOOS=windows GOARCH=amd64 go build -v .
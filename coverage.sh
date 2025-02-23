#!/bin/sh

t="/tmp/go-cover.$$.tmp"
cd .. & go test ./... -coverprofile=$t $@ && go tool cover -html=$t && unlink $t
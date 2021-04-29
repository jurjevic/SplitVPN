#!/usr/bin/env bash

if git diff-index --quiet HEAD --; then
    echo "File changes cleaned up..."
else
    echo "Changes found! Please commit dem first"
    # exit 1
fi

go install github.com/jurjevic/golf@latest
go get github.com/blang/semver/v4

$(go env GOPATH)/bin/golf -v version.go version.go -- '
'

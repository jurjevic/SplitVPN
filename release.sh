#!/usr/bin/env bash

if [ -z "$1" ]; then
  echo "Please provide a semantic version to create a release. e.g. 1.0.0"
  exit 1
fi

if git diff-index --quiet HEAD --; then
    echo "File changes cleaned up..."
else
    echo "Changes found! Please commit dem first"
    # exit 1
fi

# go install github.com/jurjevic/golf@latest
go get github.com/blang/semver/v4@latest

new_version="$1"
tag_version="v$new_version"
latest="latest"

$(go env GOPATH)/bin/golf -v version.go version.go -- '
  var NewVersion string = "'$new_version'"
'
git add version.go
git commit -m "$tag_version release build with version increment."

git push origin :refs/tags/$latest
git tag -fa $latest -m "$tag_version release build with version increment."
git tag -fa $tag_version -m "$tag_version release build with version increment."
git push origin main --tags

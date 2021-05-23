#!/usr/bin/env bash

if [ -z "$1" ]; then
  echo "Please provide a semantic version to create a release. e.g. 1.0.0"
  exit 1
fi

if git diff-index --quiet HEAD --; then
    echo "File changes cleaned up..."
else
    echo "Changes found! Please commit dem first"
    exit 1
fi

git show-ref --tags

new_version="$1"
tag_version="v$new_version"
latest="latest"

read -p "Continue with $tag_version (or ctrl-c to exit)"

golf -v version.go version.go -- '
  var NewVersion string = "'$new_version'"
'

golf -v README.md README.md -- '
  var NewVersion string = "'$new_version'"
'

git add version.go
git commit -m "$tag_version release build with version increment."

git push origin :refs/tags/$latest
git tag -fa $latest -m "$tag_version release build with version increment."
git tag -fa $tag_version -m "$tag_version release build with version increment."
git push origin main --tags --force

cd $(brew --repository jurjevic/homebrew-tap)
download="https://github.com/jurjevic/SplitVPN/archive/$tag_version.tar.gz"
wget $download
hash=$(sha256sum $tag_version.tar.gz)
rm "$tag_version.tar.gz"

golf Formula/splitvpn.rb Formula/splitvpn.rb -- '
  var HashOutput string = "'$hash'"
  var Hash string = Split(HashOutput, " ")[0]
  var Download string = "'$download'"
'

git add Formula/splitvpn.rb
git commit -m "splitvpn $tag_version added."
git push origin

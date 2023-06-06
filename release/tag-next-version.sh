#!/usr/bin/env bash

minorChangesTypes='feat|perf|refactor|revert'

function latestVersion() {
    git describe --tags --match 'v[0-9].[0-9].[0-9]' --abbrev=0
}

function aheadLatestVersion() {
    git describe --tags --match 'v[0-9].[0-9].[0-9]' --long | cut -d '-' -f2
}

function minorChangesCount() {
  git --no-pager  log --pretty="format:%s" HEAD..."$1" | grep -cE "^($minorChangesTypes){1}(\([[:alnum:]._-]+\))?(!)?:.*"
}

function patchChangesCount() {
  git --no-pager  log --pretty="format:%s" HEAD..."$1" | grep -cvE "^($minorChangesTypes){1}(\([[:alnum:]._-]+\))?(!)?:.*"
}

## ensure we have tags
git fetch --force --tags

{
rawVersion=$(latestVersion)
} || {
  git tag v0.0.0
  echo "New version v0.0.0 tag created"
  exit
}
version=${rawVersion:1}
patch=$(echo "$version" | cut -d '.' -f3)
minor=$(echo "$version" | cut -d '.' -f2)
major=$(echo "$version" | cut -d '.' -f1)

if (($(aheadLatestVersion) == 0)); then
  echo "Already at the latest version tag, not doing anything"
  exit 0
fi

nextVersion="$rawVersion"
if (($(minorChangesCount "$rawVersion") > 0)); then
  nextVersion="v$major.$((minor+1)).0"
elif (($(patchChangesCount "$rawVersion") > 0)); then
  nextVersion="v$major.$minor.$((patch+1))"
fi

git tag "$nextVersion"

echo "New version $nextVersion tag created"

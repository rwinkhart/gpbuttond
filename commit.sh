#!/bin/sh
gofmt -l -w -s .
git add -f .gitignore 15gpbuttond.start APKBUILD build-all-targets.sh commit.sh go.mod go.sum gpbuttond.go LICENSE README.md
git commit -m "$1"
git push

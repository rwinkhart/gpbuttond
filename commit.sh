#!/bin/sh
gofmt -l -w -s .
git add -f 15gpbuttond.start APKBUILD build-all-targets.sh commit.sh go.mod gpbuttond.go LICENSE README.md
git commit -m "$1"
git push

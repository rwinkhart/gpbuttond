#!/bin/sh
gofmt -l -w -s .
git add -f 15gpbuttond.start commit.sh go.mod gpbuttond.go LICENSE README.md
git commit -m "$1"
git push

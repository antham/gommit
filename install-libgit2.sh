#!/bin/sh

cd "$GOPATH"/src/github.com/antham/gommit/vendor/github.com/libgit2/git2go/ || exit

git fetch --all
git checkout -b next --track origin/next
git submodule update --init
make install

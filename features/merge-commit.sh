#!/bin/bash

cd test

# Create branch
git checkout -b test
touch file3
git add file3
git commit -F- <<EOF
feat(file) : new file 3

create a new file 3
EOF

# Checkout branch master
git checkout master

# Merge branch test
git merge --no-ff test

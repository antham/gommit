#!/bin/bash

cd test

# Update file 1 with a bad commit message
echo "test 2" > file1
git add file1
git commit -F- <<EOF
whatever
EOF

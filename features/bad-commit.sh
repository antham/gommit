#!/bin/bash

cd testing-repository || exit 1

# Update file 2 with a bad commit message
echo "test 2 test 2" > file2
git add file2
git commit -F- <<EOF
whatever
EOF

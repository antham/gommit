#!/bin/bash

if [ ! -z ${CI} ];
then
    git config --global user.name "whatever";
    git config --global user.email "whatever@example.com";
fi

# Configure name

# Init
rm -rf test > /dev/null
git init test
cd test

# Create file 1
touch file1
git add file1
git commit -F- <<EOF
feat(file) : new file 1

create a new file 1
EOF

# Create file 2
touch file2
git add file2
git commit -F- <<EOF
feat(file2) : new file 2

create a new file 2
EOF

# Update file 1
echo "test" > file1
git add file1
git commit -F- <<EOF
update(file1) : update file 1

update file 1 with a text
EOF

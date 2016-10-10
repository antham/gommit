#!/bin/bash

cd test

# Add file 4 with a long summary message
touch file4
git add file4
git commit -m "A very long summary commit greater than minimum length 50"

#!/bin/bash

cd testing-repository || exit 1

# Add file 4 with a long summary message
touch file9
git add file9
git commit -m "A very long summary commit greater than minimum length 50"

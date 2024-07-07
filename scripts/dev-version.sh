#!/bin/bash

latest_commit=$(git log -1 --format="%H")
echo "go get github.com/brownhounds/swift@${latest_commit} && go mod tidy"

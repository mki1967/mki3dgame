#!/bin/bash

# A naive Bash script for building the executable: mki3dgame 
# It is assumed that you are using bash shell and you have:
#  - Go compiler installed, and
#  - $GOPATH set properly, and
#  - the Go packagaes 'github.com/go-gl/{gl,glfw,mathgl}' and 'golang.org/x' installed

while getopts ":u" opt; do
    echo 'Updating Go packages ...';
    go get -u -v # update Go packages;
done;
echo 'Compiling mki3dgame ...'
go build # build locally the executable 'mki3dgame'

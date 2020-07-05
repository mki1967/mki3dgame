#!/bin/bash

# A naive Bash script for making snap
# It is assumed that you are using bash shell and you have:
#  - Go compiler installed, and
#  - $GOPATH set properly, and
#  - the Go packagaes 'github.com/go-gl/{gl,glfw,mathgl}' and 'golang.org/x' installed
#  - sapcraft installed
#  - rsync installed

while getopts ":u" opt; do
    echo 'Updating Go packages ...';
    go get -u -v # update Go packages;
done;

./make-mki3game.bash

echo 'Preparing for snap ...'

rm -rf build-snap # start with fresh directory
mkdir build-snap # do everything in the 'build' directory


pushd build-snap/ # Go to the build directory.
echo  'copy the relevant files for snap ...'
cp ../mki3dgame-snap-wrapper.bash .
rsync -av ../snap .
rsync -av ../assets .
cp ../mki3dgame .
echo 'Building snap ...'
snapcraft --use-lxd
# mv mki3dgame-snap_*.snap ../
popd    # return to the directory

echo 'Your snap should be in: ./build-snap/mki3dgame-snap_[...].snap'

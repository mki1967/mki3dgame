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

echo 'Preparing for AppImage ...'
mkdir build-snap # do everything in the 'build' directory


#  It is assumed that this directory is called 'mki3dgame'
pushd ../ # Go to the parent directory.
echo 'rsync-ing ...'
rsync -av --exclude-from=mki3dgame/rsync-exclude-patterns-snap mki3dgame mki3dgame/build-snap/
popd    # return to the directory

pushd build-snap/mki3dgame # Go to the build directory.
echo 'Building snap ...'
snapcraft
mv mki3dgame-snap_*.snap ../
cd ../
rm -rf mki3dgame/ # remove the rsync-ed directory
popd    # return to the directory

echo 'Your snap should be in: ./build-snap/mki3dgame-snap_[...].snap'

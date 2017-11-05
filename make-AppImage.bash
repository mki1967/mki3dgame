#!/bin/bash

# A naive Bash script for making AppImage
# It is assumed that you are using bash shell and you have:
#  - Go compiler installed, and
#  - $GOPATH set properly, and
#  - the Go packagaes 'github.com/go-gl/{gl,glfw,mathgl}' and 'golang.org/x' installed
#  - the appimagetool-x86_64.AppImage on your $PATH (see https://github.com/AppImage/AppImageKit )
#  - rsync installed

while getopts ":u" opt; do
    echo 'Updating Go packages ...';
    go get -u -v # update Go packages;
done;


./make-mki3game.bash

echo 'Preparing for AppImage ...'
mkdir build-AppImage # do everything in the 'build' directory


#  It is assumed that this directory is called 'mki3dgame'
pushd ../ # Go to the parent directory.
echo 'rsync-ing ...'
rsync -av --exclude-from=mki3dgame/rsync-exclude-patterns-AppImage mki3dgame mki3dgame/build-AppImage/
popd    # return to the directory

pushd build-AppImage/ # Go to the build directory.
echo 'Building AppImage ...'
appimagetool-x86_64.AppImage mki3dgame
rm -rf mki3dgame/ # remove the rsync-ed directory
popd    # return to the directory

echo 'Your AppImage should be in: ./build-AppImage/mki3dgame-x86_64.AppImage'

#!/bin/bash

# A naive Bash script for making AppImage

echo 'Compiling ...'
go build # build locally the executable 'mki3dgame'

mkdir build-AppImage # do everything in the 'build' directory


#  It is assumed that this directory is called 'mki3dgame'
pushd ../ # Go to the parent directory.
echo 'rsync-ing ...'
rsync -avz --exclude-from=mki3dgame/rsync-exclude-patterns-AppImage mki3dgame mki3dgame/build-AppImage/
popd    # return to the directory

pushd build-AppImage/ # Go to the build directory.
echo 'Building AppImage ...'
appimagetool-x86_64.AppImage mki3dgame
rm -rf mki3dgame/ # remove the rsync-ed directory
popd    # return to the directory

echo 'Your AppImage should be in: ./build-AppImage/mki3dgame-x86_64.AppImage'

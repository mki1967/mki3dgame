#!/bin/sh
export LIBGL_DRIVERS_PATH=$SNAP/usr/lib/x86_64-linux-gnu/dri
# export MKI3DGAME_ASSETS=$SNAP/assets
exec $SNAP/mki3dgame $SNAP/assets

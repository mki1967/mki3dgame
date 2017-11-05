#!/bin/sh
export LIBGL_DRIVERS_PATH=$SNAP/usr/lib/x86_64-linux-gnu/dri
# export MKI3DGAME_ASSETS=$SNAP/assets
export XDG_DATA_DIRS="$SNAP/usr/share:$XDG_DATA_DIRS"
exec $SNAP/mki3dgame $SNAP/assets

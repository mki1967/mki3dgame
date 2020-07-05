#!/bin/sh
export LIBGL_DRIVERS_PATH=$SNAP/usr/lib/x86_64-linux-gnu/dri
# export MKI3DGAME_ASSETS=$SNAP/assets
export XDG_DATA_DIRS="$SNAP/usr/share:$XDG_DATA_DIRS"
echo
echo '**Gamepad configuration command:**   `snap connect mki3dgame-snap:joystick`'
echo
exec $SNAP/mki3dgame $SNAP/assets

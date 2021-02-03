#!/bin/bash
unset WINDOWID
zenity --text-info --font='DejaVu Sans Mono' --width 700 --height 700 --title "MKI3D GAME HELP PAGE" <<EOF
MKI3D GAME (see: https://github.com/mki1967/mki3dgame )
=======================================================

Collect the tokens and avoid the monsters.
(You can be captured by a monster ;-) 


PRESS THE MOUSE BUTTON ON THE SCREEN SECTORS

Action-sectors layout:
----------------------
+---------+---------+---------+
| MF      |   MU    | MF      |
|         +---------+         |
|         |   RU    |         |
+---------+---------+---------+
|    |    |         |    |    |
| ML | RL |   LV    | RR | MR |
|    |    |         |    |    |
+---------+---------+---------+
|         |   RD    |         |
|         +---------+         |
| MB      |   MD    | MB      |
+---------+---------+---------+

MF, MB, MU, MD, ML, MR - move forward, backward, up, down, left and right
RU, RD, RL, RR - rotate up, down, left and right
LV - set the player upright, then align the horizontal rotation to the right angle

OR USE THE KEYBOARD

Key press actions:
------------------

  H - this help message
  Arrow keys - rotate observer
  Shift + Arrow keys  - move observer sideways
  Space - set the player upright and then align with the horizontal axes
  B, F - move observer backward or forward
  L - set the diffuse light direction perpendicular to the screen
  P - pause the game
  X - load next random stage
  Q - display your score and the number of remaining tokens
  S - toggle SKYBOX ON/OFF
  N - new SKYBOX
  Shift+Q - quit
  F11 - toggle full-screen

OR USE THE GAMEPAD:

Gamepad actions:
----------------

  Left Joystick  - rotate
  Button A   - set the player upright and then align with the horizontal axes
  Right Joystick - move forward/backward/left/right
  Right/Left Triger - forward/backward
  Dpad Buttons - move on the screen plane 
  Button Start -  load next random stage

EOF


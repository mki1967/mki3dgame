package main

const helpText = `
MKI3DGAME (see: https://github.com/mki1967/mki3dgame )
======================================================

Collect the tokens and avoid the monsters.
(You can be captured by a monster ;-) 

Use the keyboard:

Key press actions:
------------------

  H - this help message
  Arrow keys - rotate observer
  Shift + Arrow keys  - move observer sideways
  Space - set the player upright
  B, F - move observer backward or forward
  L - set the diffuse light direction perpendicular to the screen
  P - pause the game
  X - reload random stage

Or press the mouse on the screen sectors:

Sectors action layout:
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
LV - set the player upright
`

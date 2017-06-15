package main

const helpText = `
Key press actions:
------------------

  H - this help screen
  Arrow keys - rotate observer
  Shift + Arrow keys  - move observer sideways
  B, F - move observer backward or forward
  L - set the diffuse light direction perpendicular to the screen
  P - pause the game
  X - reload random stage

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
`

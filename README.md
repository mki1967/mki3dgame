# mki3dgame [![Build Status](https://travis-ci.org/mki1967/mki3dgame.svg?branch=master)](https://travis-ci.org/mki1967/mki3dgame)

This game uses designs produced with MKI3D web 3D editor (see: https://mki1967.github.io/mki3d/ ).

See also on: 
[![Download mki3dgame](https://sourceforge.net/sflogo.php?type=15&group_id=2849958)](https://sourceforge.net/p/mki3dgame/)
or Bintray: 
[ ![Download](https://api.bintray.com/packages/mki1967/mki3d-generic/mki3dgame/images/download.svg) ](https://bintray.com/mki1967/mki3d-generic/mki3dgame/_latestVersion)

Snapped version is available at: https://uappexplorer.com/snap/ubuntu/mki3dgame-snap


In the game, you have to collect tokens scattered in the stages and avoid being captured by the monsters.
(A short screen-cast is available at: https://youtu.be/vp6nhvOqhdU . )
Run the game with the path to assets directory as the command line argument.
(See the content of the 'runme' script in this directory.)

The assets directory has the following sub-directories:

* 'icons' -  icon '.png' files (some systems may use them ...)
* 'monsters' - monster shapes '.mki3d' files - made with MKI3D
* 'sectors'  - shapes of screen sectors '.mki3d' - made with MKI3D, specific to the code 
* 'stages'  - stages '.mki3d' files - made with MKI3D
* 'tokens'  - token shapes '.mki3d' files - made with MKI3D

You can design your own stages and the shapes of monsters or tokens
with this editor.
Just place the files in the respective sub-directories
'stages', 'monsters', or 'tokens' of the main assets directory.
Shapes are selected randomly from each sub-directory for each stage.

To build the the game from the source code with Go compiler you need the following packages:
*	"github.com/go-gl/gl/v3.3-core/gl"
*	"github.com/go-gl/glfw/v3.2/glfw"
*	"github.com/go-gl/mathgl/mgl32"
*	"github.com/mki1967/go-mki3d/mki3d"
*	"github.com/mki1967/go-mki3d/glmki3d"

This project has been moved here from the collection of Go program demos at https://github.com/mki1967/test-go-mki3d.git


INSTALLATION FROM THE SOURCE CODE REPOSITORY
--------------------------------------------

* Install Go language on your system, unless you already have it. (See: https://golang.org/doc/install , do not forget to set the `GOPATH` environment variable)
* Install requirements for your system for gl and glfw listed in the README.md files:
    - https://github.com/go-gl/gl/blob/master/README.md
    - https://github.com/go-gl/glfw/blob/master/README.md
    
* Try to use the command: `go get -u github.com/mki1967/mki3dgame` (It can take some time to complete ...)
* The binary compiled file should be installed in `${GOPATH}/bin` directory.
* Note that you can provide your own path to the assets directory as a command line argument.
  Otherwise the program will try to find the `assets` directory either in  `${GOPATH}/src/github.com/mki1967/mki3dgame/`
  or (if `GOPATH` is not set) in the current directory.
* You can also prepare your own AppImage for Linux as described in https://github.com/mki1967/mki3dgame/blob/master/appImage-instructions.txt

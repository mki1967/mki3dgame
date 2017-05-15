// run in the source code directory with: go run *.go  <filename>.mki3d

package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/go-gl/mathgl/mgl32"
	// "github.com/mki1967/go-mki3d/glmki3d"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	rand.Seed(time.Now().Unix()) // init random generator
}

const windowWidth = 800
const windowHeight = 600

var Window *glfw.Window // main window

func message(msg string) error {
	fmt.Println(msg)
	err := Window.Iconify()
	if err != nil {
		return err
	}
	fmt.Print("(PRESS ENTER TO RESUME:)")
	fmt.Scanln()
	err = Window.Restore()
	if err != nil {
		panic(err)
	}
	// err = Window.Maximize()
	Window.Show()
	fmt.Println("RESUMED.")
	return err
}

var doInMainThread func() = nil

func main() {

	// get path to assets from command line argument
	if len(os.Args) < 2 {
		panic(errors.New(" *** PROVIDE PATH TO ASSETS DIRECTORY AS A COMMAND LINE ARGUMENT !!! *** "))
	}

	// fragments from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4) // try multisampling for better quality ...
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	Window = window // copy to global variable
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// version := gl.GoStr(gl.GetString(gl.VERSION))
	// fmt.Println("OpenGL version", version)

	// Configure global settings
	gl.Enable(gl.MULTISAMPLE) // probably not needed ...
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.3, 1.0)

	// message(helpText) // initial help message
	fmt.Println(helpText)

	// init game structure from assets and the game window
	game, err := MakeEmptyGame(os.Args[1], window)
	if err != nil {
		panic(err)
	}

	// load the first stage
	err = game.Init()
	if err != nil {
		panic(err)
	}

	// main loop
	for !window.ShouldClose() {

		// Maintenance
		window.SwapBuffers()
		if doInMainThread != nil {
			doInMainThread()     // execute required function
			doInMainThread = nil // done
		}

		// if( game.Paused ) { // old version
		if game.Paused.Get() { // new version
			game.CancelAction()
			glfw.WaitEvents()
			game.ProbeTime()
		} else {
			game.ProbeTime()
			game.Update()
			// game.Paused = game.PauseRequest.TestAndCancel() // check for pause
			// game.Paused.Set(game.PauseRequest.TestAndCancel()) // check for pause -- new version
			glfw.PollEvents()
		}
		game.Redraw()

	}

	fmt.Println("YOUR TOTAL SCORE IS: ", game.TotalScore)
	// cleanup
	game.StageDSPtr.DeleteData()
	game.FrameDSPtr.DeleteData()
	game.SectorsDSPtr.DeleteData()
	game.TokenDSPtr.DeleteData()
	game.MonsterDSPtr.DeleteData()

}

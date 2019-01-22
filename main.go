// run in the source code directory with: go run *.go  <filename>.mki3d

package main

import (
	// "errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"strconv"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/go-gl/mathgl/mgl32"
	// "github.com/mki1967/go-mki3d/glmki3d"
	"log"
	"math/rand"
	"os"
	// "path/filepath"
	"runtime"
	// "strings"
	"flag"
	"time"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	rand.Seed(time.Now().Unix()) // init random generator
}

var Window *glfw.Window // main window

func message(msg string) error {
	fmt.Println(msg)
	err := error(nil)
	/* May cause problems if the process runs in the background ...
	err = Window.Iconify()
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
	*/
	return err
}

var doInMainThread func() = nil

var PathToAssets string = "assets"

func pathToAssets() string {
	pathToAssets := "assets"
	/*
		gopath, isGopath := os.LookupEnv("GOPATH")
		// pathToAssets := os.Getenv("GOPATH") + "/src/github.com/mki1967/mki3dgame/assets"
		if isGopath {
			pathToAssets = gopath + "/src/github.com/mki1967/mki3dgame/assets"
			fmt.Println("If you have built from source code, then  you should have some assests in: " + pathToAssets)
		}

		// check if it is installed in '/usr/...' directory  -- requires Go version >= 1.8
		execPath, err := os.Executable()
		fmt.Printf("execPath = %v\n", execPath) // test
		if err == nil && strings.Contains(execPath, "/usr/") {
			execDir := filepath.Dir(execPath)
			// relative '../share/games/mki3game/assets/' to
			pathToAssets = filepath.Dir(execDir) + "/share/games/mki3game/assets/"
			fmt.Println("If you have installed from distribution, then  you should have some assests in: " + pathToAssets)
		}
	*/
	envPath, isEnvPath := os.LookupEnv("MKI3DGAME_ASSETS")
	if isEnvPath {
		// fmt.Println("environment variable MKI3DGAME_ASSETS is set to " + envPath) // test
		pathToAssets = envPath
	}

	// get path to assets from command line argument
	// if len(os.Args) < 2 {
	if len(flag.Args()) < 1 {
		// panic(errors.New(" *** PROVIDE PATH TO ASSETS DIRECTORY AS A COMMAND LINE ARGUMENT !!! *** "))
		// fmt.Println(" *** YOU CAN PROVIDE PATH TO YOUR ASSETS DIRECTORY AS A COMMAND LINE ARGUMENT !!! *** ")
		fmt.Println(" Trying to use default path to assets: '" + pathToAssets + "'")

	} else {
		// fmt.Println("Path to assets from command line argument: " + os.Args[1])
		// pathToAssets = os.Args[1]
		pathToAssets = flag.Args()[0]
	}

	return pathToAssets
}

func main() {

	windowWidthPtr := flag.Int("width", 800, "initial width of the window")
	windowHeightPtr := flag.Int("height", 600, "initial height of the window")
	flag.Parse()

	// fmt.Println(flag.Args())

	pathToAssets := pathToAssets()
	PathToAssets = pathToAssets

	// test for zenity
	ZenityTest()

	// fragments from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	// glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4) // try multisampling for better quality ...
	window, err := glfw.CreateWindow(*windowWidthPtr, *windowHeightPtr, "Mki3dgame", nil, nil)
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
	gl.DepthFunc(gl.LEQUAL)
	gl.ClearColor(0.0, 0.0, 0.3, 1.0)

	// init game structure from assets and the game window
	game, err := MakeEmptyGame(pathToAssets, window)
	if err != nil {
		panic(err)
	}

	// load the first stage
	err = game.Init()
	if err != nil {
		panic(err)
	}

	// message(helpText) // initial help message

	firstLine := "MKI3D GAME with " + strconv.Itoa(len(game.AssetsPtr.Stages)) + " stages.\n"

	fmt.Println(firstLine + helpText)
	ZenityInfo(firstLine+"\nPRESS THE MOUSE ON SCREEN SECTORS OR USE THE KEYS.\n\nPRESS 'H' FOR HELP. ", "6")
	// doInMainThread = ZenityHelp
	game.Paused.Set(true)
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
	ts := strconv.FormatFloat(game.TotalScore, 'f', 2, 64)
	ZenityInfo("YOUR TOTAL SCORE IS: "+ts, "4")
	// cleanup
	game.StageDSPtr.DeleteData()
	game.FrameDSPtr.DeleteData()
	game.SectorsDSPtr.DeleteData()
	game.TokenDSPtr.DeleteData()
	game.MonsterDSPtr.DeleteData()

}

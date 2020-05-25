package main

import (
	"fmt"
	// "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	// "github.com/go-gl/mathgl/mgl32"
	// "math"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/mki1967/go-mki3d/glmki3d"
)

func (g *Mki3dGame) CheckGamepad() {

	if glfw.Joystick1.Present() && glfw.Joystick1.IsGamepad() {
		fmt.Println("Joystick1.Present && Joystick1.IsGamepad") /// test

		gamepadStatePtr := glfw.Joystick1.GetGamepadState()

		fmt.Println("*gamepadStatePtr = %v", *gamepadStatePtr) /// test

	}
	/*
		g.JustCollected = false // stop celebrations
		if action == glfw.Release {
			g.CancelAction()
			return
		}

		if action == glfw.Press {
			// g.Paused = false // old version
			g.Paused.Set(false) // new version
			width, height := w.GetSize()
			fx, fy := w.GetCursorPos()
			x := int(fx)
			y := int(fy)
			g.SetSectorAction(x, y, width, height)
			return
		}
	*/

}

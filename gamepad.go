package main

import (
	// "fmt"
	// "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	// "github.com/go-gl/mathgl/mgl32"
	// "math"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/mki1967/go-mki3d/glmki3d"
)

func (g *Mki3dGame) CheckGamepad() {

	if glfw.Joystick1.Present() && glfw.Joystick1.IsGamepad() {
		// fmt.Println("Joystick1.Present && Joystick1.IsGamepad") /// test

		gamepadStatePtr := glfw.Joystick1.GetGamepadState()

		// fmt.Println("*gamepadStatePtr = %v", *gamepadStatePtr) /// test

		var nextAction ActionIndex = ActionNIL

		switch {
		case gamepadStatePtr.Axes[glfw.AxisRightX] < -0.5:
			nextAction = ActionRL
		case gamepadStatePtr.Axes[glfw.AxisRightX] > 0.5:
			nextAction = ActionRR

		case gamepadStatePtr.Axes[glfw.AxisRightY] < -0.5:
			nextAction = ActionRU
		case gamepadStatePtr.Axes[glfw.AxisRightY] > 0.5:
			nextAction = ActionRD

		case gamepadStatePtr.Axes[glfw.AxisLeftY] < -0.5:
			nextAction = ActionMU
		case gamepadStatePtr.Axes[glfw.AxisLeftY] > 0.5:
			nextAction = ActionMD

		case gamepadStatePtr.Axes[glfw.AxisLeftX] < -0.5:
			nextAction = ActionML
		case gamepadStatePtr.Axes[glfw.AxisLeftX] > 0.5:
			nextAction = ActionMR

		case gamepadStatePtr.Axes[glfw.AxisRightTrigger] > -0.75:
			nextAction = ActionMF
		case gamepadStatePtr.Axes[glfw.AxisLeftTrigger] > -0.75:
			nextAction = ActionMB

		case gamepadStatePtr.Buttons[glfw.ButtonA] == glfw.Press:
			nextAction = ActionLV
		}

		if nextAction == g.LastGamepadAction {
			return // continuation of the same action or noaction
		}

		if nextAction != g.LastGamepadAction {
			g.CancelAction() // the action is changed -- reset speed
		}

		g.LastGamepadAction = nextAction       // record  as the last gamepad action
		g.SetAction(g.ActionArray[nextAction]) // set current action

	}

}

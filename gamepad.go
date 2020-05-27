package main

import (
	"fmt"
	// "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	// "github.com/go-gl/mathgl/mgl32"
	// "math"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/mki1967/go-mki3d/glmki3d"
	// "time"
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
			nextAction = ActionRD
		case gamepadStatePtr.Axes[glfw.AxisRightY] > 0.5:
			nextAction = ActionRU

		case gamepadStatePtr.Axes[glfw.AxisLeftY] < -0.5:
			nextAction = ActionMF
		case gamepadStatePtr.Axes[glfw.AxisLeftY] > 0.5:
			nextAction = ActionMB

		case gamepadStatePtr.Axes[glfw.AxisLeftX] < -0.5:
			nextAction = ActionML
		case gamepadStatePtr.Axes[glfw.AxisLeftX] > 0.5:
			nextAction = ActionMR

		case gamepadStatePtr.Axes[glfw.AxisRightTrigger] > -0.95:
			nextAction = ActionMF
		case gamepadStatePtr.Axes[glfw.AxisLeftTrigger] > -0.95:
			nextAction = ActionMB

		case gamepadStatePtr.Buttons[glfw.ButtonDpadUp] == glfw.Press:
			nextAction = ActionMU
		case gamepadStatePtr.Buttons[glfw.ButtonDpadDown] == glfw.Press:
			nextAction = ActionMD

		case gamepadStatePtr.Buttons[glfw.ButtonDpadLeft] == glfw.Press:
			nextAction = ActionML
		case gamepadStatePtr.Buttons[glfw.ButtonDpadRight] == glfw.Press:
			nextAction = ActionMR

		case gamepadStatePtr.Buttons[glfw.ButtonA] == glfw.Press:
			nextAction = ActionLV

		case gamepadStatePtr.Buttons[glfw.ButtonStart] == glfw.Press:
			fmt.Println("RELOADING RANDOM STAGE ...")
			ZenityInfo("NEXT RANDOM STAGE ...", "1")
			g.NextStage()

		}

		if nextAction != ActionNIL {
			g.Paused = false // new version for single-thread version
			g.LastActivityTime = glfw.GetTime()
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

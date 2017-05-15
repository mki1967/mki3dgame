package main

import (
	"fmt"
	// "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	// "math"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/mki1967/go-mki3d/glmki3d"
)

func (g *Mki3dGame) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release {
		if action == glfw.Release {
			g.CancelAction()
			return
		}

	}

	if action == glfw.Press {
		// g.Paused = false // old version
		g.Paused.set(false) // new version

		// g := GamePtr // short name
		switch {

		/* rotate player */
		case key == glfw.KeyRight && mods == 0:
			g.SetAction(g.ActionRotRight)
		case key == glfw.KeyLeft && mods == 0:
			g.SetAction(g.ActionRotLeft)

		case key == glfw.KeyUp && mods == 0:
			g.SetAction(g.ActionRotUp)

		case key == glfw.KeyDown && mods == 0:
			g.SetAction(g.ActionRotDown)

		case key == glfw.KeySpace:
			g.SetAction(g.ActionLevel)

			/* move player */
		case key == glfw.KeyRight && mods == glfw.ModShift:
			g.SetAction(g.ActionMoveRight)

		case key == glfw.KeyLeft && mods == glfw.ModShift:
			g.SetAction(g.ActionMoveLeft)

		case key == glfw.KeyUp && mods == glfw.ModShift:
			g.SetAction(g.ActionMoveUp)

		case key == glfw.KeyDown && mods == glfw.ModShift:
			g.SetAction(g.ActionMoveDown)

		case key == glfw.KeyF && mods == glfw.ModShift:
			fallthrough
		case key == glfw.KeyF && mods == 0, key == glfw.KeyEnter:
			g.SetAction(g.ActionMoveForward)

		case key == glfw.KeyB && mods == glfw.ModShift:
			fallthrough
		case key == glfw.KeyB && mods == 0, key == glfw.KeyBackspace:
			g.SetAction(g.ActionMoveBackward)

		case key == glfw.KeyL && mods == 0: /* light */
			g.StageDSPtr.UniPtr.LightUni = g.StageDSPtr.UniPtr.ViewUni.Mat3().Inv().Mul3x1(mgl32.Vec3{0, 0, 1}).Normalize()

		case key == glfw.KeyZ && mods == 0: /* zoom out */
			width, height := w.GetSize()
			zy := g.StageDSPtr.Mki3dPtr.Projection.ZoomY / 1.1
			fmt.Println("ZoomY: ", zy)
			g.StageDSPtr.Mki3dPtr.Projection.ZoomY = zy
			g.StageDSPtr.UniPtr.SetProjectionFromMki3d(g.StageDSPtr.Mki3dPtr, width, height)
		case key == glfw.KeyZ && mods == glfw.ModShift: /* zoom in */
			width, height := w.GetSize()
			zy := g.StageDSPtr.Mki3dPtr.Projection.ZoomY * 1.1
			fmt.Println("ZoomY: ", zy)
			g.StageDSPtr.Mki3dPtr.Projection.ZoomY = zy
			g.StageDSPtr.UniPtr.SetProjectionFromMki3d(g.StageDSPtr.Mki3dPtr, width, height)

			/* help */
		case key == glfw.KeyH && mods == 0:
			message(helpText)
		case key == glfw.KeyP && mods == 0: /* PAUSE */
			g.PauseRequest.Set()
			fmt.Println("PAUSED")
		}
	}
}

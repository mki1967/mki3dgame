package main

import (
	// "fmt"
	// "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	// "github.com/go-gl/mathgl/mgl32"
	// "math"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/mki1967/go-mki3d/glmki3d"
)

func (g *Mki3dGame) Mki3dMouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Release {
		g.CancelAction()
		return
	}

	if action == glfw.Press {
		g.Paused = false
		width, height := w.GetSize()
		fx, fy := w.GetCursorPos()
		x := int(fx)
		y := int(fy)
		g.SetSectorAction(x, y, width, height)
		return
	}

}

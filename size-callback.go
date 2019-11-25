package main

import (
	// "fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	// "github.com/go-gl/mathgl/mgl32"
	// "math"
	// "github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/go-mki3d/glmki3d"
)

// Function to be used as resize callback
func (g *Mki3dGame) SizeCallback(w *glfw.Window, width int, height int) {
	// g := GamePtr                                                                                                 // short name
	gl.Viewport(0, 0, int32(width), int32(height))                                                                // inform GL about resize
	g.StageDSPtr.UniPtr.ProjectionUni = glmki3d.ProjectionMatrix(g.StageDSPtr.Mki3dPtr.Projection, width, height) // recompute projection matrix
	// fmt.Println("SizeCallback ",  width, " ", height)
}

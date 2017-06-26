package glmki3d

import (
	// "errors"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/mki3d"
)

// GLUni - values of parameters to be stored in shaders' uniforms
type GLUni struct {
	ProjectionUni mgl32.Mat4
	ViewUni       mgl32.Mat4
	ModelUni      mgl32.Mat4
	LightUni      mgl32.Vec3
	AmbientUni    float32
}

// SetSimple  sets GLUni for simple drawing directly in clipping space
// (i.e. orthogonal projection projection on the XY plane in the cube [-1,1]^3 resized to the window)
func (glUni *GLUni) SetSimple() {
	glUni.ProjectionUni = mgl32.Ident4()
	glUni.ViewUni = mgl32.Ident4()
	glUni.ModelUni = mgl32.Ident4()
	glUni.LightUni = mgl32.Vec3{0, 0, 1}
	glUni.AmbientUni = 1

}

// MakeGLUni makes GLUni with values of uniforms set to simple default values  and returns a pointer to it.
func MakeGLUni() *GLUni {
	var glUni GLUni

	glUni.SetSimple()

	return &glUni
}

// ProjectionMatrix computes GL projection matrix from mki3d.ProjectionType, using width and height of current window
func ProjectionMatrix(p mki3d.ProjectionType, width, height int) mgl32.Mat4 {
	// make both width and height greater than zero
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	h := float32(height)
	w := float32(width)
	xx := p.ZoomY * h / w
	yy := p.ZoomY
	zz := (p.ZFar + p.ZNear) / (p.ZFar - p.ZNear)
	zw := float32(1.0)
	wz := -2 * p.ZFar * p.ZNear / (p.ZFar - p.ZNear)

	var m mgl32.Mat4

	m.SetRow(0, mgl32.Vec4{xx, 0, 0, 0})
	m.SetRow(1, mgl32.Vec4{0, yy, 0, 0})
	m.SetRow(2, mgl32.Vec4{0, 0, zz, wz})
	m.SetRow(3, mgl32.Vec4{0, 0, zw, 0})

	return m

}

// Mat3 converts Matrix3dType to mgl32.Mat3
func Mat3(m mki3d.Matrix3dType) mgl32.Mat3 {
	var q mgl32.Mat3
	q.SetRow(0, mgl32.Vec3(m[0]))
	q.SetRow(1, mgl32.Vec3(m[1]))
	q.SetRow(2, mgl32.Vec3(m[2]))
	return q
}

// ViewMatrix computes GL view matrix from mki3d.ViewType
func ViewMatrix(v mki3d.ViewType) mgl32.Mat4 {

	mov := mgl32.Vec3(v.FocusPoint).Mul(-1)

	rot := Mat3(v.RotationMatrix).Mul(v.Scale)
	scrSh := v.ScreenShift

	m := rot.Mat4()
	m.SetCol(3, mgl32.Vec4{mov.Dot(rot.Row(0)) + scrSh[0], mov.Dot(rot.Row(1)) + scrSh[1], mov.Dot(rot.Row(2)) + scrSh[2], 1.0})

	//
	return m

}

func (glUni *GLUni) SetModelPosition(pos mgl32.Vec3) {
	glUni.ModelUni.SetCol(3, pos.Vec4(1))
}

// Sets uniform projection matrix of glUni based on the data from mki3dData and on the width and height (of the display window).
// The function panics if mki3dData == nil
func (glUni *GLUni) SetProjectionFromMki3d(mki3dData *mki3d.Mki3dType, width, height int) {
	glUni.ProjectionUni = ProjectionMatrix(mki3dData.Projection, width, height)
}

// Sets the view matrix of glUni based on the data from mki3dData.
// The function panics if mki3dData == nil
func (glUni *GLUni) SetViewFromMki3d(mki3dData *mki3d.Mki3dType) {
	glUni.ViewUni = ViewMatrix(mki3dData.View)
}

// Sets the light parameters of glUni based on the data from mki3dData
// The function panics if mki3dData == nil
func (glUni *GLUni) SetLightFromMki3d(mki3dData *mki3d.Mki3dType) {
	glUni.LightUni = mgl32.Vec3(mki3dData.Light.Vector)
	glUni.AmbientUni = mki3dData.Light.AmbientFraction
}

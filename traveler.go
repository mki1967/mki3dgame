package main

import (
	// "fmt" // tests
	// "errors"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const TravelerMovSpeed = 15 // units per second
const TravelerRotSpeed = 90 // degrees per second

type Traveler struct {
	Position   mgl32.Vec3 // position
	Rot        RotHVType  // orientation
	MovSpeed   float64
	RotSpeed   float64
	CapturedBy *MonsterType
}

func (t *Traveler) ViewMatrix() mgl32.Mat4 {
	c1 := float32(math.Cos(-t.Rot.XZ * degToRadians))
	s1 := float32(math.Sin(-t.Rot.XZ * degToRadians))

	c2 := float32(math.Cos(-t.Rot.YZ * degToRadians))
	s2 := float32(math.Sin(-t.Rot.YZ * degToRadians))

	v := t.Rot.ViewerRotatedVector(t.Position.Mul(-1))

	// row-major ??
	return mgl32.Mat4{
		c1, 0, -s1, v[0],
		-s2 * s1, c2, -s2 * c1, v[1],
		c2 * s1, s2, c2 * c1, v[2],
		0, 0, 0, 1,
	}.Transpose()
}

func (t *Traveler) Move(dx, dy, dz float32) {
	v := t.Rot.WorldRotatedVector(mgl32.Vec3{dx, dy, dz})
	t.Position = t.Position.Add(v)

}

func (t *Traveler) ClipToBox(vmin, vmax mgl32.Vec3) {

	if t.Position[0] > vmax[0] {
		t.Position[0] = vmax[0]
	}
	if t.Position[0] < vmin[0] {
		t.Position[0] = vmin[0]
	}

	if t.Position[1] > vmax[1] {
		t.Position[1] = vmax[1]
	}
	if t.Position[1] < vmin[1] {
		t.Position[1] = vmin[1]
	}

	if t.Position[2] > vmax[2] {
		t.Position[2] = vmax[2]
	}
	if t.Position[2] < vmin[2] {
		t.Position[2] = vmin[2]
	}

}

func MakeTraveler(position mgl32.Vec3) *Traveler {
	var t Traveler
	t.Position = position
	t.MovSpeed = TravelerMovSpeed
	t.RotSpeed = TravelerRotSpeed
	return &t
}

func (t *Traveler) Update(g *Mki3dGame) {
	if t.CapturedBy != nil {
		t.Position = t.CapturedBy.Position
		g.StageDSPtr.UniPtr.ViewUni = g.TravelerPtr.ViewMatrix()
	}
}

package main

import (
	// "fmt" // tests
	// "errors"
	"github.com/go-gl/mathgl/mgl32"
	// "github.com/mki1967/go-mki3d/mki3d"
	// "github.com/mki1967/go-mki3d/glmki3d"
	"math"
	"math/rand"
)

// Random position in the game stage box with the margin offset from the borders
func (game *Mki3dGame) RandPosition(margin float32) mgl32.Vec3 {
	m := mgl32.Vec3{margin, margin, margin}
	v1 := game.VMin.Add(m)
	v2 := game.VMax.Sub(m)
	return RandPosition(v1, v2)
}

// Random position in the box [vmin, vmax]
func RandPosition(vmin, vmax mgl32.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		rand.Float32()*(vmax[0]-vmin[0]) + vmin[0],
		rand.Float32()*(vmax[1]-vmin[1]) + vmin[1],
		rand.Float32()*(vmax[2]-vmin[2]) + vmin[2],
	}
}

// RotHVType represents sequence of two rotations:
// by the angle XY on XY-plane and by the angle YZ on YZ-plane
// (in degrees)
type RotHVType struct {
	XZ float64
	YZ float64
}

// Rotate XZ by angle (in degrees)
func (rot *RotHVType) RotateXZ(angle float64) {
	rot.XZ = math.Remainder(rot.XZ+angle, 360)
}

// Rotate YZ by angle (in degrees) an clamping result to [-90,90]
func (rot *RotHVType) RotateYZ(angle float64) {
	rot.YZ = math.Remainder(rot.YZ+angle, 360)
	if rot.YZ < -90 {
		rot.YZ = -90
	}
	if rot.YZ > 90 {
		rot.YZ = 90
	}
}

// find the nearest right angle (in dergrees) among the angles {0, 90, 180, 270}
func NearestRightAngle(angle float64) float64 {
	angle = angle - math.Floor(angle/360)*360
	d := math.Abs(angle)
	out := float64(0)

	x := math.Abs(angle - 90)
	if x < d {
		out = float64(90)
		d = x
	}

	x = math.Abs(angle - 180)
	if x < d {
		out = float64(180)
		d = x
	}

	x = math.Abs(angle - 270)

	if x < d {
		out = float64(270)
		d = x
	}

	x = math.Abs(angle - 360)

	if x < d {
		out = float64(0)
		d = x
	}

	return out

}

const degToRadians = math.Pi / 180

func (rot *RotHVType) WorldRotatedVector(vector mgl32.Vec3) mgl32.Vec3 {
	c1 := float32(math.Cos(rot.XZ * degToRadians))
	s1 := float32(math.Sin(rot.XZ * degToRadians))
	c2 := float32(math.Cos(rot.YZ * degToRadians))
	s2 := float32(math.Sin(rot.YZ * degToRadians))

	return mgl32.Vec3{
		c1*vector[0] - s1*s2*vector[1] - s1*c2*vector[2],
		c2*vector[1] - s2*vector[2],
		s1*vector[0] + c1*s2*vector[1] + c1*c2*vector[2],
	}
}

func (rot *RotHVType) ViewerRotatedVector(vector mgl32.Vec3) mgl32.Vec3 {
	c1 := float32(math.Cos(-rot.XZ * degToRadians))
	s1 := float32(math.Sin(-rot.XZ * degToRadians))
	c2 := float32(math.Cos(-rot.YZ * degToRadians))
	s2 := float32(math.Sin(-rot.YZ * degToRadians))

	return mgl32.Vec3{
		c1*vector[0] - s1*vector[2],
		-s2*s1*vector[0] + c2*vector[1] - s2*c1*vector[2],
		c2*s1*vector[0] + s2*vector[1] + c2*c1*vector[2],
	}
}

func RandRotated(vec mgl32.Vec3) mgl32.Vec3 {
	var rot RotHVType
	rot.XZ = rand.Float64() * 360
	rot.YZ = rand.Float64() * 360
	return rot.WorldRotatedVector(vec)
}

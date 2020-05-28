package main

import (
// "fmt" // tests
// "errors"
// "github.com/go-gl/mathgl/mgl32"
// "github.com/mki1967/go-mki3d/mki3d"
// "github.com/mki1967/go-mki3d/glmki3d"
// "math"
// "math/rand"
)

type ActionIndex uint8

const (
	ActionNIL = ActionIndex(iota)

	ActionMF
	ActionMB
	ActionMU
	ActionMD
	ActionML
	ActionMR

	ActionRU
	ActionRD
	ActionRL
	ActionRR

	ActionLV
)

const NumberOfActions = ActionLV + 1

func (g *Mki3dGame) InitActionSectors() {
	mf := func() {
		g.ActionMoveForward()
	}
	mb := func() {
		g.ActionMoveBackward()
	}

	mu := func() {
		g.ActionMoveUp()
	}
	md := func() {
		g.ActionMoveDown()
	}

	ml := func() {
		g.ActionMoveLeft()
	}
	mr := func() {
		g.ActionMoveRight()
	}

	rl := func() {
		g.ActionRotLeft()
	}
	rr := func() {
		g.ActionRotRight()
	}

	ru := func() {
		g.ActionRotUp()
	}
	rd := func() {
		g.ActionRotDown()
	}

	lv := func() {
		g.ActionLevel()
	}

	g.ActionArray = [12]func(){ // the same sequence as in ActionIndex
		nil,
		mf, mb, mu, md, ml, mr, // moves
		ru, rd, rl, rr, // rotations
		lv, // level
	}
	g.ActionSectors = [6][6]func(){
		{mf, mf, mu, mu, mf, mf},
		{mf, mf, ru, ru, mf, mf},
		{ml, rl, lv, lv, rr, mr},
		{ml, rl, lv, lv, rr, mr},
		{mb, mb, rd, rd, mb, mb},
		{mb, mb, md, md, mb, mb},
	}

}

func (g *Mki3dGame) SetSectorAction(x, y, width, height int) {
	sx := HorizontalSectors * x / width
	sy := VerticalSectors * y / height
	g.SetAction(g.ActionSectors[sy][sx])

}

func (g *Mki3dGame) SetAction(action func()) {
	g.CurrentAction = action
}

func (g *Mki3dGame) CancelAction() {
	g.CurrentAction = nil
	g.TravelerPtr.ResetSpeed()
	g.JustCollected = false // stop celebrations
}

func (g *Mki3dGame) ActionMoveForward() {
	d := float32(g.TravelerPtr.MovSpeed * g.LastTimeDelta)
	g.TravelerPtr.Move(0, 0, d)
	g.TravelerPtr.ClipToBox(g.VMin, g.VMax)
	g.TravelerPtr.UpdateSpeed(g.LastTimeDelta)
}

func (g *Mki3dGame) ActionMoveBackward() {
	d := float32(g.TravelerPtr.MovSpeed * g.LastTimeDelta)
	g.TravelerPtr.Move(0, 0, -d)
	g.TravelerPtr.ClipToBox(g.VMin, g.VMax)
	g.TravelerPtr.UpdateSpeed(g.LastTimeDelta)
}

func (g *Mki3dGame) ActionMoveLeft() {
	d := float32(g.TravelerPtr.MovSpeed * g.LastTimeDelta)
	g.TravelerPtr.Move(-d, 0, 0)
	g.TravelerPtr.ClipToBox(g.VMin, g.VMax)
	g.TravelerPtr.UpdateSpeed(g.LastTimeDelta)
}

func (g *Mki3dGame) ActionMoveRight() {
	d := float32(g.TravelerPtr.MovSpeed * g.LastTimeDelta)
	g.TravelerPtr.Move(d, 0, 0)
	g.TravelerPtr.ClipToBox(g.VMin, g.VMax)
	g.TravelerPtr.UpdateSpeed(g.LastTimeDelta)
}

func (g *Mki3dGame) ActionMoveUp() {
	d := float32(g.TravelerPtr.MovSpeed * g.LastTimeDelta)
	g.TravelerPtr.Move(0, d, 0)
	g.TravelerPtr.ClipToBox(g.VMin, g.VMax)
	g.TravelerPtr.UpdateSpeed(g.LastTimeDelta)
}

func (g *Mki3dGame) ActionMoveDown() {
	d := float32(g.TravelerPtr.MovSpeed * g.LastTimeDelta)
	g.TravelerPtr.Move(0, -d, 0)
	g.TravelerPtr.ClipToBox(g.VMin, g.VMax)
	g.TravelerPtr.UpdateSpeed(g.LastTimeDelta)
}

func (g *Mki3dGame) ActionRotLeft() {
	d := g.TravelerPtr.RotSpeed * g.LastTimeDelta
	g.TravelerPtr.Rot.RotateXZ(d)
}

func (g *Mki3dGame) ActionRotRight() {
	d := g.TravelerPtr.RotSpeed * g.LastTimeDelta
	g.TravelerPtr.Rot.RotateXZ(-d)
}

func (g *Mki3dGame) ActionRotUp() {
	d := g.TravelerPtr.RotSpeed * g.LastTimeDelta
	g.TravelerPtr.Rot.RotateYZ(-d)
}

func (g *Mki3dGame) ActionRotDown() {
	d := g.TravelerPtr.RotSpeed * g.LastTimeDelta
	g.TravelerPtr.Rot.RotateYZ(d)
}

func (g *Mki3dGame) ActionLevel() {
	if g.TravelerPtr.Rot.YZ == 0 {
		// fmt.Println("1: g.TravelerPtr.Rot.XZ ==", g.TravelerPtr.Rot.XZ)
		g.TravelerPtr.Rot.XZ = NearestRightAngle(g.TravelerPtr.Rot.XZ)
		// fmt.Println("1: g.TravelerPtr.Rot.XZ ==", g.TravelerPtr.Rot.XZ)
	} else {
		g.TravelerPtr.Rot.YZ = 0
	}
	g.CancelAction()
}

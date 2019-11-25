package main

import (
	"fmt" // tests
	// "strconv"
	// "errors"
	"github.com/go-gl/mathgl/mgl32"
	// "github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/go-mki3d/glmki3d"
	// "math"
	// "math/rand"
)

// Parameters of a single token
type TokenType struct {
	Position  mgl32.Vec3
	Collected bool
	DSPtr     *glmki3d.DataShader // shape for redraw (may be shared by many)
}

// Creates a token  at position pos with datashader *dsptr
func MakeToken(pos mgl32.Vec3, dsPtr *glmki3d.DataShader) *TokenType {
	var t TokenType
	t.Position = pos
	t.DSPtr = dsPtr
	t.Collected = false
	return &t
}

// Redraw token m
func (t *TokenType) Draw() {
	if t.Collected {
		return
	}

	t.DSPtr.UniPtr.SetModelPosition(t.Position)
	t.DSPtr.DrawModel()
}

// square of the distance to collect token
const TokenCollectionSqrDist = 1

// Update token t in game g
func (t *TokenType) Update(g *Mki3dGame) {
	if t.Collected {
		return // should not be considered
	}
	v := t.Position.Sub(g.TravelerPtr.Position)
	if v.Dot(v) < TokenCollectionSqrDist {
		t.Collected = true
		g.TokensRemaining--
		g.TokensCollected++
		g.CancelAction()
		// some celebrations ...
		fmt.Println("COLLECTED !!! (", g.TokensRemaining, " remaining)")
		// ZenityInfo("COLLECTED !!! ("+ strconv.Itoa(g.TokensRemaining)+" remaining)","1")
		g.JustCollected = true

	}
}

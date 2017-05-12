package main

import (
	"fmt" // tests
	// "errors"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/go-mki3d/glmki3d"
	"math"
	// "math/rand"
	_ "image/png"
)

const BoxMargin = 30 // margin for bounding box of the stage

var FrameColor = mki3d.Vector3dType{1.0, 1.0, 1.0} // color of the bounding box frame

var NumberOfMonsters = 20

var NumberOfTokens = 10

const VerticalSectors = 6   // vertical dimmension of sectors array
const HorizontalSectors = 6 // horizontal  dimmension of sectors array

// data structure for the game
type Mki3dGame struct {
	// assets info
	AssetsPtr *Assets
	// GLFW data
	WindowPtr *glfw.Window
	// GL shaders
	ShaderPtr *glmki3d.Shader
	// Shape data shaders
	StageDSPtr   *glmki3d.DataShader
	FrameDSPtr   *glmki3d.DataShader // frame of the bounding box (computed for the stage)
	SectorsDSPtr *glmki3d.DataShader
	TokenDSPtr   *glmki3d.DataShader
	MonsterDSPtr *glmki3d.DataShader

	VMin, VMax mgl32.Vec3 // corners of the bounding box of the stage (computed with the BoxMargin)

	TravelerPtr *Traveler // the first person (the player)

	Monsters []*MonsterType // set of monsters

	Tokens []*TokenType // set of tokens

	TokensRemaining int     // number of remaining tokens
	TokensCollected int     // number of remaining tokens
	TotalScore      float64 // number of remaining tokens

	StageStartingTime float64 // game global time probing
	LastProbedTime    float64 // game global time probing
	LastTimeDelta     float64 // game global time probing

	CurrentAction func()                                     // current action of the player
	ActionSectors [VerticalSectors][HorizontalSectors]func() // functions of the mouse actions

	PauseRequest Flag // set by a goroutine to request pause
	Paused bool // true if game is paused

	WasAction Flag // set ech time the user action is executed
}

// Make game structure with the shader and without any data.
// Prepare assets info using pathToAssets.
// Return pointer to the strucure.
// To be used once (before the game).
func MakeEmptyGame(pathToAssets string, window *glfw.Window) (*Mki3dGame, error) {
	var game Mki3dGame

	game.WindowPtr = window

	shaderPtr, err := glmki3d.MakeShader()
	if err != nil {
		return nil, err
	}

	game.ShaderPtr = shaderPtr

	assetsPtr, err := LoadAssets(pathToAssets)
	if err != nil {
		return nil, err
	}
	game.AssetsPtr = assetsPtr

	game.InitActionSectors()

	// set icons
	imgs, err := assetsPtr.LoadIcons()
	if err != nil {
		return nil, err
	}

	window.SetIcon(imgs)

	// setting the callbacks
	window.SetSizeCallback(game.SizeCallback)
	window.SetKeyCallback(game.KeyCallback)
	window.SetMouseButtonCallback(game.Mki3dMouseButtonCallback)

	game.PauseRequest = MakeFlag()
	game.WasAction = MakeFlag()

	go game.EcoFreezer() // run concurrent eco-freezer goroutine
	
	return &game, nil

}

// Load data and init game for each new stage.
func (game *Mki3dGame) Init() (err error) {

	width, height := game.WindowPtr.GetSize()

	err = game.InitSectors()
	if err != nil {
		return err
	}

	err = game.InitStage(width, height)
	if err != nil {
		return err
	}

	err = game.InitToken()
	if err != nil {
		return err
	}

	err = game.InitMonster()
	if err != nil {
		return err
	}

	fmt.Println("NEW STAGE! Collect ", game.TokensRemaining, " tokens.")
	// init time probe
	game.LastProbedTime = glfw.GetTime()
	game.StageStartingTime = game.LastProbedTime
	game.LastTimeDelta = 0

	return nil
}

// Use this function to  before the actions.
func (game *Mki3dGame) ProbeTime() {
	now := glfw.GetTime()
	game.LastTimeDelta = now - game.LastProbedTime
	game.LastProbedTime = now
}

// Load sectors shape and init the SectorsDSPtr.
func (game *Mki3dGame) InitSectors() error {

	sectorsPtr, err := game.AssetsPtr.LoadRandomSectors()
	if err != nil {
		return err
	}

	sectorsDataShaderPtr, err := glmki3d.MakeDataShader(game.ShaderPtr, sectorsPtr)
	if err != nil {
		return err
	}

	sectorsDataShaderPtr.UniPtr.SetSimple()

	if game.SectorsDSPtr != nil {
		game.SectorsDSPtr.DeleteData() // free old GL buffers
	}

	game.SectorsDSPtr = sectorsDataShaderPtr

	return nil
}

// Load token shape and init the tokenDSPtr.
func (game *Mki3dGame) InitToken() error {

	tokenPtr, err := game.AssetsPtr.LoadRandomToken()
	if err != nil {
		return err
	}

	tokenDataShaderPtr, err := glmki3d.MakeDataShader(game.ShaderPtr, tokenPtr)
	if err != nil {
		return err
	}

	tokenDataShaderPtr.UniPtr.SetSimple()

	if game.TokenDSPtr != nil {
		game.TokenDSPtr.DeleteData() // free old GL buffers
	}

	game.TokenDSPtr = tokenDataShaderPtr
	game.GenerateTokens()

	return nil
}

func (game *Mki3dGame) GenerateTokens() {
	game.Tokens = make([]*TokenType, NumberOfTokens)
	for i := range game.Tokens {
		game.Tokens[i] = MakeToken(game.RandPosition(BoxMargin), game.TokenDSPtr)
	}
	game.TokensRemaining = NumberOfTokens
}

func (game *Mki3dGame) DrawTokens() {
	for _, t := range game.Tokens {
		t.Draw()
	}
}

func (game *Mki3dGame) UpdateTokens() {
	for _, t := range game.Tokens {
		t.Update(game)
	}
}

// Load monster shape and init the monsters.
func (game *Mki3dGame) InitMonster() error {

	monsterPtr, err := game.AssetsPtr.LoadRandomMonster()
	if err != nil {
		return err
	}

	monsterDataShaderPtr, err := glmki3d.MakeDataShader(game.ShaderPtr, monsterPtr)
	if err != nil {
		return err
	}

	monsterDataShaderPtr.UniPtr.SetSimple()

	if game.MonsterDSPtr != nil {
		game.MonsterDSPtr.DeleteData() // free old GL buffers
	}

	game.MonsterDSPtr = monsterDataShaderPtr

	// game.MonsterDSPtr.UniPtr.SetModelPosition(game.RandPosition(BoxMargin)) // test
	game.GenerateMonsters()

	return nil
}

func (game *Mki3dGame) GenerateMonsters() {
	game.Monsters = make([]*MonsterType, NumberOfMonsters)
	for i := range game.Monsters {
		game.Monsters[i] = MakeMonster(game.RandPosition(0), game.MonsterDSPtr)
	}
}

func (game *Mki3dGame) DrawMonsters() {
	for _, m := range game.Monsters {
		m.Draw()
	}
}

func (game *Mki3dGame) UpdateMonsters() {
	for _, m := range game.Monsters {
		m.Update(game)
	}
}

const InitStageZoomY = 2.5

// Load stage shape and init the related data.
func (game *Mki3dGame) InitStage(width, height int) error {

	stagePtr, err := game.AssetsPtr.LoadRandomStage()
	if err != nil {
		return err
	}

	stageDataShaderPtr, err := glmki3d.MakeDataShader(game.ShaderPtr, stagePtr)
	if err != nil {
		return err
	}

	stageDataShaderPtr.UniPtr.SetSimple()
	stageDataShaderPtr.Mki3dPtr.Projection.ZoomY = InitStageZoomY
	stageDataShaderPtr.UniPtr.SetProjectionFromMki3d(stageDataShaderPtr.Mki3dPtr, width, height)
	// stageDataShaderPtr.UniPtr.SetProjectionFromMki3d(stagePtr, width, height)
	stageDataShaderPtr.UniPtr.SetLightFromMki3d(stagePtr)

	stageDataShaderPtr.UniPtr.ViewUni = mgl32.Ident4()
	stageDataShaderPtr.UniPtr.ViewUni.SetCol(3, mgl32.Vec3(stageDataShaderPtr.Mki3dPtr.Cursor.Position).Mul(-1).Vec4(1))

	if game.StageDSPtr != nil {
		game.StageDSPtr.DeleteData() // free old GL buffers
	}

	game.StageDSPtr = stageDataShaderPtr

	game.copmuteVMinVMax() // compute bounding box of the stage: VMin, VMax
	game.copmuteFrame()    // visible line frame of the bounding box

	game.TravelerPtr = MakeTraveler(mgl32.Vec3(stagePtr.Cursor.Position))

	return nil
}

// recompute bounding box  with the BoxMargin  corners of the stage.
func (game *Mki3dGame) copmuteVMinVMax() {
	stagePtr := game.StageDSPtr.Mki3dPtr
	game.VMax = mgl32.Vec3(stagePtr.Cursor.Position) // cursror position should be included - the starting poin of traveler
	game.VMin = game.VMax

	for _, seg := range stagePtr.Model.Segments {
		for _, point := range seg {
			for d := range point.Position {
				if game.VMax[d] < point.Position[d] {
					game.VMax[d] = point.Position[d]
				}
				if game.VMin[d] > point.Position[d] {
					game.VMin[d] = point.Position[d]
				}
			}

		}
	}

	for _, tr := range stagePtr.Model.Triangles {
		for _, point := range tr {
			for d := range point.Position {
				if game.VMax[d] < point.Position[d] {
					game.VMax[d] = point.Position[d]
				}
				if game.VMin[d] > point.Position[d] {
					game.VMin[d] = point.Position[d]
				}
			}

		}
	}

	m := mgl32.Vec3{BoxMargin, BoxMargin, BoxMargin}

	game.VMin = game.VMin.Sub(m)
	game.VMax = game.VMax.Add(m)
}

// recompute frame of the bounding box corners of the stage.
func (game *Mki3dGame) copmuteFrame() {
	a := game.VMin
	b := game.VMax

	v000 := mki3d.Vector3dType(a)
	v001 := mki3d.Vector3dType{a[0], a[1], b[2]}
	v010 := mki3d.Vector3dType{a[0], b[1], a[2]}
	v011 := mki3d.Vector3dType{a[0], b[1], b[2]}
	v100 := mki3d.Vector3dType{b[0], a[1], a[2]}
	v101 := mki3d.Vector3dType{b[0], a[1], b[2]}
	v110 := mki3d.Vector3dType{b[0], b[1], a[2]}
	v111 := mki3d.Vector3dType(b)

	lines := [][2]mki3d.Vector3dType{
		{v000, v001},
		{v010, v011},
		{v100, v101},
		{v110, v111},

		{v000, v010},
		{v001, v011},
		{v100, v110},
		{v101, v111},

		{v000, v100},
		{v001, v101},
		{v010, v110},
		{v011, v111}}

	segments := mki3d.SegmentsType(make([]mki3d.SegmentType, 12))

	for i := range segments {
		segments[i] = mki3d.SegmentType{
			{Position: lines[i][0], Color: FrameColor},
			{Position: lines[i][1], Color: FrameColor}}
	}

	var frameMki3d mki3d.Mki3dType

	frameMki3d.Model.Segments = segments

	dsPtr, err := glmki3d.MakeDataShader(game.ShaderPtr, &frameMki3d)

	if err != nil {
		panic(err)
	}

	dsPtr.UniPtr.SetSimple()

	if game.FrameDSPtr != nil {
		game.FrameDSPtr.DeleteData() // free old GL buffers
	}

	game.FrameDSPtr = dsPtr

}

func (game *Mki3dGame) Update() {
	game.UpdateMonsters()
	game.UpdateTokens()
	game.TravelerPtr.Update(game) // Captured ?
	// check the state
	if game.TokensRemaining <= 0 { // go to next stage
		// compute some results ...
		time := math.Floor(game.LastProbedTime - game.StageStartingTime)
		score := math.Floor(1 + 30*float64(game.TokensCollected)/(time+1))
		game.TotalScore += score
		fmt.Println("STAGE FINISHED !!! Time:", time, " seconds, stage score: ", score, ", total: ", game.TotalScore, ".")
		game.NextStage()
	} else { // update the player in the stage
		if game.CurrentAction != nil {
			game.CurrentAction()
			game.StageDSPtr.UniPtr.ViewUni = game.TravelerPtr.ViewMatrix()
		}

	}
}

func (game *Mki3dGame) NextStage() {
	// block callbacks ... ?
	// reload next stage
	game.TokensCollected = 0
	game.Init()
	// unblock callbacks ... ?
}

// Redraw the game stage
func (game *Mki3dGame) Redraw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // to be moved to redraw ?
	// draw stage
	game.StageDSPtr.SetBackgroundColor()
	game.StageDSPtr.DrawStage()
	// draw frame
	game.FrameDSPtr.DrawModel()
	// draw tokens
	game.DrawTokens()
	// draw monsters
	game.DrawMonsters()

	if game.CurrentAction == nil {
		// draw sectors
		gl.Disable(gl.DEPTH_TEST)
		game.SectorsDSPtr.DrawStage()
		gl.Enable(gl.DEPTH_TEST)
	} else {
		game.WasAction.Set(); // set for EcoFreezer
	}

}

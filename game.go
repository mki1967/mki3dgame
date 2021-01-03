package main

import (
	"fmt" // tests
	// "errors"
	"math"
	"strconv"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/glmki3d"
	"github.com/mki1967/go-mki3d/mki3d"
	"github.com/mki1967/go-skybox/sbxgpu"

	// "time"
	_ "image/png"
	"math/rand"
)

const BoxMargin = 30 // margin for bounding box of the stage

var FrameColor = mki3d.Vector3dType{1.0, 1.0, 1.0} // color of the bounding box frame

var NumberOfMonsters = 20

var NumberOfTokens = 10

var TokenInfoPositions = squareSpiral(NumberOfTokens)

const VerticalSectors = 6   // vertical dimmension of sectors array
const HorizontalSectors = 6 // horizontal  dimmension of sectors array

const skyboxChance = 0.25 // the chance of having withSkybox==true on the new stage

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

	LastActivityTime float64 // recorded time of last activity for the auto-pause

	CurrentAction func()                                     // current action of the player
	ActionSectors [VerticalSectors][HorizontalSectors]func() // functions of the mouse actions
	ActionArray   [NumberOfActions]func()                    // indexed functions of the actions

	LastGamepadAction ActionIndex // last action invoked by the gamepad

	// PauseRequest Flag // set by a goroutine to request pause
	Paused bool // true if game is paused -- used again for the single-thread version
	// Paused SharedBool // true if game is paused - shared version

	// WasAction Flag // set ech time the user action is executed -- not used in the single-thread version

	JustCollected bool // token has been just collected! Do some clebrations in Redraw ...

	Skybox     sbxgpu.SbxGpu // skybox
	withSkybox bool          // draw the skybox?
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

	// game.PauseRequest = MakeFlag()
	// game.WasAction = MakeFlag() // -- not used in the single-thread version

	// game.Paused = MakeSharedBool() // new version -- not used in the single-thread version

	game.Skybox = sbxgpu.NewSbxGpu() // init skybox shader

	// go game.EcoFreezer() // run concurrent eco-freezer goroutine -- no goroutines in the single-thread version

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

	err = game.InitToken(width, height)
	if err != nil {
		return err
	}

	err = game.InitMonster()
	if err != nil {
		return err
	}

	game.Skybox.RenderRandomCube() // make random cube
	game.withSkybox = (rand.Float64() < skyboxChance)
	fmt.Println("NEW STAGE! Collect ", game.TokensRemaining, " tokens.")
	// ZenityInfo("NEW STAGE! Collect "+strconv.Itoa( game.TokensRemaining)+ " tokens.", "2")
	// init time probe
	game.LastProbedTime = glfw.GetTime()
	game.StageStartingTime = game.LastProbedTime
	game.LastTimeDelta = 0

	game.Paused = true // start in paused state

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
func (game *Mki3dGame) InitToken(width int, height int) error {

	tokenPtr, err := game.AssetsPtr.LoadRandomToken()
	if err != nil {
		return err
	}

	tokenDataShaderPtr, err := glmki3d.MakeDataShader(game.ShaderPtr, tokenPtr)
	if err != nil {
		return err
	}

	// tokenDataShaderPtr.UniPtr.SetSimple()
	tokenDataShaderPtr.UniPtr.SetProjectionFromMki3d(tokenDataShaderPtr.Mki3dPtr, width, height)
	tokenDataShaderPtr.UniPtr.SetViewFromMki3d(tokenDataShaderPtr.Mki3dPtr)
	tokenDataShaderPtr.UniPtr.SetLightFromMki3d(tokenDataShaderPtr.Mki3dPtr)

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
		t := strconv.FormatFloat(time, 'f', 2, 64)
		s := strconv.FormatFloat(score, 'f', 2, 64)
		ts := strconv.FormatFloat(game.TotalScore, 'f', 2, 64)
		roundInfo := ""
		if game.AssetsPtr.LastLoadedStage+1 == len(game.AssetsPtr.Stages) {
			roundInfo = "\nROUND FINISHED !!!"
		}
		fmt.Println("STAGE FINISHED !!! Time:", time, " seconds, stage score: ", score, ", total: ", game.TotalScore, ".")
		ZenityInfo("STAGE FINISHED !!! Time:"+t+" seconds,\n stage score: "+s+",\n total: "+ts+"."+roundInfo, "5")
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

	if game.JustCollected {
		gl.ClearColor(0.0, 0.4, 0.4, 1.0)
		game.TokenDSPtr.UniPtr.SetModelPosition(mgl32.Vec3{0, 0, 0})
		game.TokenDSPtr.DrawStage()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		for i := 0; i < game.TokensRemaining; i++ {
			game.TokenDSPtr.UniPtr.SetModelPosition(TokenInfoPositions[i])
			game.TokenDSPtr.DrawModel()
		}
		game.Skybox.RenderRandomCube()
		game.withSkybox = true
		// game.JustCollected = false
		// time.Sleep(time.Millisecond * 500)
	} else {
		// draw stage
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // to be moved to redraw ?
		game.StageDSPtr.SetBackgroundColor()
		game.StageDSPtr.DrawStage()
		// draw frame
		game.FrameDSPtr.DrawModel()
		// draw tokens
		game.DrawTokens()
		// draw monsters
		game.DrawMonsters()
		if game.withSkybox {
			game.Skybox.DrawSkybox(game.StageDSPtr.UniPtr.ViewUni, game.StageDSPtr.UniPtr.ProjectionUni) // draw the skybox
		}

	}
	if game.CurrentAction == nil && !game.JustCollected {
		// draw sectors
		gl.Disable(gl.DEPTH_TEST)
		game.SectorsDSPtr.DrawStage()
		gl.Enable(gl.DEPTH_TEST)
		game.tryAutoPause()
	} else {
		game.LastActivityTime = glfw.GetTime()
		// game.WasAction.Set() // set for EcoFreezer -- not used in the single-thread version
	}

}

const autoPauseTimeDelta = 15 // seconds since last activity for auto-pause

func (game *Mki3dGame) tryAutoPause() {
	if game.Paused {
		return // do not pause if already paused
	}

	if glfw.GetTime()-game.LastActivityTime > autoPauseTimeDelta {
		game.Paused = true
		fmt.Println("AUTO-PAUSE")
		ZenityInfo("AUTO-PAUSE", "1")

	}
}

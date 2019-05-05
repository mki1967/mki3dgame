package main

import (
	"fmt" // tests
	"io/ioutil"
	// "errors"
	"github.com/mki1967/go-mki3d/mki3d"
	"image"
	_ "image/png"
	"math/rand"
	"os"
)

// the structure for assets of mki3dgame
type Assets struct {
	Path            string
	Assets          []os.FileInfo
	Stages          []os.FileInfo
	Tokens          []os.FileInfo
	Sectors         []os.FileInfo
	Monsters        []os.FileInfo
	Icons           []os.FileInfo
	LastLoadedStage int // should be initialised to -1
	/// TestStages [10000]bool // for tests if we have less tham 10000 stages ;-)
}

/*
func testStages( a *Assets ) {
	for i:= range a.Stages {
		if ! a.TestStages[ i ] {
	                fmt.Println("i==", i," len(a.Stages)==", len(a.Stages))
			return
		}
	}
	fmt.Println("ALL STAGES HAVE BEEN LOADED !!!")
}
*/

const (
	StagesDir   = "stages"
	TokensDir   = "tokens"
	SectorsDir  = "sectors"
	MonstersDir = "monsters"
	IconsDir    = "icons"
	PS          = string(os.PathSeparator)
)

func LoadAssets(pathToAssets string) (*Assets, error) {
	ass, err := ioutil.ReadDir(pathToAssets)
	if err != nil {
		return nil, err
	}

	assets := Assets{Path: pathToAssets, Assets: ass, LastLoadedStage: -1} /// ...

	assets.Stages, err = ioutil.ReadDir(pathToAssets +
		PS +
		StagesDir)

	if err != nil {
		return &assets, err
	}

	assets.Tokens, err = ioutil.ReadDir(pathToAssets +
		PS +
		TokensDir)

	if err != nil {
		return &assets, err
	}

	assets.Monsters, err = ioutil.ReadDir(pathToAssets +
		PS +
		MonstersDir)

	if err != nil {
		return &assets, err
	}

	assets.Sectors, err = ioutil.ReadDir(pathToAssets +
		PS +
		SectorsDir)

	if err != nil {
		return &assets, err
	}

	assets.Icons, err = ioutil.ReadDir(pathToAssets +
		PS +
		IconsDir)

	if err != nil {
		return &assets, err
	}

	return &assets, nil
}

func (a *Assets) load(dir string, fname string) (*mki3d.Mki3dType, error) {
	mki3dPtr, err := mki3d.ReadFile(a.Path +
		PS +
		dir +
		PS +
		fname)
	if err != nil {
		return nil, err
	}
	return mki3dPtr, err
}

func (a *Assets) LoadIcons() ([]image.Image, error) {
	img := make([]image.Image, len(a.Icons))
	for i, f := range a.Icons {
		name := a.Path + PS + IconsDir + PS + f.Name()
		imgFile, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		img[i], _, err = image.Decode(imgFile)
		if err != nil {
			return nil, err
		}
		// fmt.Println("Loaded: ", name) // test

	}
	return img, nil
}

/* old version
func (a *Assets) LoadRandomStage() (*mki3d.Mki3dType, error) {
	stages := len(a.Stages)
	n := stages

	if stages >= 2 && a.LastLoadedStage >= 0 { // if we have at least 2 stages and at something has been loaded
		n = stages - 1
	}

	// r := rand.Intn(len(a.Stages))  // old
	r := (a.LastLoadedStage + 1 + rand.Intn(n)) % stages // if n==stages-1 then should select any stage different from  the last one
	mki3dPtr, err := a.load(StagesDir, a.Stages[r].Name())
	a.LastLoadedStage = r // record the index of the loaded stage
	/// a.TestStages[r]= true /// for tests
	/// testStages(a)
	return mki3dPtr, err
}
*/

func (a *Assets) randomShuffle(t *[]os.FileInfo) {
	n := len(*t)

	if n <= 1 { // nothing to shuffle
		return
	}

	swap := func(i, j int) {
		tmp := (*t)[i]
		(*t)[i] = (*t)[j]
		(*t)[j] = tmp
	}

	for i := 0; i < n; i++ {
		j := i + rand.Intn(n-i)
		swap(i, j)
	}

}

func (a *Assets) LoadRandomStage() (*mki3d.Mki3dType, error) {

	stages := len(a.Stages)
	r := (a.LastLoadedStage + 1) % stages
	if r == 0 {
		fmt.Println("SHUFFLING ...")
		a.randomShuffle(&(a.Stages)) // reshuffle for next round
		for i, s := range a.Stages {
			fmt.Println(i, s.Name())
		}
	}
	mki3dPtr, err := a.load(StagesDir, a.Stages[r].Name())
	fmt.Println(a.Stages[r].Name(), r)
	a.LastLoadedStage = r // record the index of the loaded stage
	return mki3dPtr, err
}

func (a *Assets) LoadRandomToken() (*mki3d.Mki3dType, error) {
	r := rand.Intn(len(a.Tokens))
	mki3dPtr, err := a.load(TokensDir, a.Tokens[r].Name())
	return mki3dPtr, err
}

func (a *Assets) LoadRandomMonster() (*mki3d.Mki3dType, error) {
	r := rand.Intn(len(a.Monsters))
	mki3dPtr, err := a.load(MonstersDir, a.Monsters[r].Name())
	return mki3dPtr, err
}

func (a *Assets) LoadRandomSectors() (*mki3d.Mki3dType, error) {
	r := rand.Intn(len(a.Sectors))
	mki3dPtr, err := a.load(SectorsDir, a.Sectors[r].Name())
	return mki3dPtr, err
}

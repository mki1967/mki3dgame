package main

import (
	// "fmt" // tests
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
	Path     string
	Assets   []os.FileInfo
	Stages   []os.FileInfo
	Tokens   []os.FileInfo
	Sectors  []os.FileInfo
	Monsters []os.FileInfo
	Icons    []os.FileInfo
}

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

	assets := Assets{Path: pathToAssets, Assets: ass} /// ...

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

func (a *Assets) LoadRandomStage() (*mki3d.Mki3dType, error) {
	r := rand.Intn(len(a.Stages))
	mki3dPtr, err := a.load(StagesDir, a.Stages[r].Name())
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

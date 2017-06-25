package mki3d

/* json operations */

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// Reads all from the input with JSON representation of MKI3d data
// Returns pointer to Mki3dType or nil and error.
func ReadAll(r io.Reader) (*Mki3dType, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var dat Mki3dType

	if err := json.Unmarshal(data, &dat); err != nil {
		panic(err)
	}

	if err != nil {
		return nil, err
	}
	return &dat, nil
}

// Reads all from the text file with JSON representation of MKI3d data.
// Returns pointer to Mki3dType or nil and error.
func ReadFile(filename string) (*Mki3dType, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	mki3dPtr, err := ReadAll(f)
	if err != nil {
		return nil, err
	}
	return mki3dPtr, nil
}

// Returns string with JSON representation of  Mki3dType
func Stringify(data *Mki3dType) (jsonOut string) {
	mki3dBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	jsonOut = string(mki3dBytes)
	return
}

package main

import (
	"fmt"
	"os/exec"
)

// Assume that zenity is available
var IsZenity = true

func ZenityTest() {
	out, err := RunScript("zenity-test.bash")
	if err != nil {
		IsZenity = false
		fmt.Println("zenity disabled")
	}

	fmt.Println("zenity version: " + string(out))
}

func RunScript(name string, arg ...string) ([]byte, error) {
	// fmt.Println(arg)
	cmd := exec.Command(PathToAssets+"/scripts/"+name, arg...)
	out, err := cmd.CombinedOutput()
	// fmt.Println(out)
	// fmt.Println(err)

	return out, err
}

func ZenityHelp() {
	if !IsZenity {
		return
	}
	go RunScript("zenity-help.bash")
}

func ZenityInfo(info string, timeout string) {
	if !IsZenity {
		return
	}
	go RunScript("zenity-info.bash", info, timeout)
}

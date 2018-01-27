package main

import (
	"fmt"
	"os/exec"
)

func RunScript(name string) {
	cmd := exec.Command(PathToAssets + "/scripts/" + name)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
	}
	/*
		fmt.Println(string(out))
	*/
}

func ZenityHelp() {
	go RunScript("zenity-help.bash")
}

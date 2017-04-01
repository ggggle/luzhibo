//+build windows

package main

import (
	"fmt"
	"github.com/pkg/browser"
	"runtime"
	"os/exec"
)

func cmd() {
	fmt.Scanln()
}

func openWebUI(hta bool) {
	u := "http://localhost"
	if port != 80 {
		u = fmt.Sprintf("%s:%d", u, port)
	}
	if runtime.GOOS == "windows" && hta {
		err := exec.Command("mshta", u+"/hta").Start()
		if err != nil {
			browser.OpenURL(u)
		}
	} else {
		browser.OpenURL(u)
	}
}

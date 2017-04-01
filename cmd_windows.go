//+build windows

package main

import (
	"fmt"
	"github.com/pkg/browser"
	"runtime"
	"os/exec"
	"syscall"
	"unsafe"
)

func cmd() {
	fmt.Scanln()
}

func openWebUI(hta bool) {
	u := "http://localhost"
	if port != 80 {
		u = fmt.Sprintf("%s:%d", u, port)
	}
	if runtime.GOOS == "windows" && hta && checkWin10() {
		err := exec.Command("mshta", u+"/hta").Start()
		if err != nil {
			browser.OpenURL(u)
		}
	} else {
		browser.OpenURL(u)
	}
}

func checkWin10() (ret bool) {
	defer func() {
		if recover() != nil {
			ret = false
		}
	}()
	type OSVERSIONINFOW struct {
		dwOSVersionInfoSize uint32
		dwMajorVersion      uint32
		dwMinorVersion      uint32
		dwBuildNumber       uint32
		dwPlatformId        uint32
		szCSDVersion        [128]byte
	}
	var v OSVERSIONINFOW
	p := uintptr(unsafe.Pointer(&v))
	kernel32 := syscall.NewLazyDLL("ntdll.dll")
	c := kernel32.NewProc("RtlGetVersion")
	r, _, _ := c.Call(p)
	if r == 0 {
		ret = v.dwMajorVersion == 10
	} else {
		ret = false
	}
	return
}

//go:generate goversioninfo -icon=icon.ico

package main

import (
	"fmt"

	"runtime"
)

const ver = 2017033000
const p = "录直播"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	s := ":12216"
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	go startServer(s)
	openWebUI()
	cmd()
}

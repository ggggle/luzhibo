//go:generate goversioninfo -icon=icon.ico -manifest luzhibo.manifest

package main

import (
	"runtime"
	"flag"
	"strconv"
	"fmt"
	"time"
	"path/filepath"
	"os"
)

const ver = 2017040300
const p = "录直播"

var port = 12216

var nhta *bool

var htaproc *os.Process

func main() {
	p := flag.Int("port", port, "WebUI监听端口")
	nopen := flag.Bool("nopenui", false, "不自动打开WebUI")
	nhta = flag.Bool("nhta", false, "禁用hta(仅Windows有效)")
	pid := flag.Int("pid", 0, "pid")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if *pid != 0 {
		proc, err := os.FindProcess(*pid)
		if err == nil {
			proc.Kill()
			time.Sleep(time.Second * 2)
		}
	}
	go func() {
		time.Sleep(time.Second * 5)
		d, f := filepath.Split(os.Args[0])
		tp := filepath.Join(d, "."+f+".old")
		os.Remove(tp)
	}()
	port = *p
	s := ":" + strconv.Itoa(port)
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	go startServer(s)
	if !*nopen {
		openWebUI(!*nhta)
	}
	cmd()
}

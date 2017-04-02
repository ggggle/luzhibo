//go:generate goversioninfo -icon=icon.ico -manifest luzhibo.manifest

package main

import (
	"runtime"
	"flag"
	"os"
	"path"
	"strconv"
	"fmt"
	"time"
)

const ver = 2017040100
const p = "录直播"

var port = 12216

var nhta *bool

func main() {
	p := flag.Int("port", port, "WebUI监听端口")
	nopen := flag.Bool("nopenui", false, "不自动打开WebUI")
	nhta = flag.Bool("nhta", false, "禁用hta(仅Windows有效)")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		time.Sleep(time.Second * 5)
		d, f := path.Split(os.Args[0])
		tp := path.Join(d, "."+f+".old")
		if _, err := os.Stat(tp); err == nil {
			os.Remove(tp)
		}
		port = *p
	}()
	s := ":" + strconv.Itoa(port)
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	go startServer(s)
	if !*nopen {
		openWebUI(!*nhta)
	}
	cmd()
}

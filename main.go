//go:generate goversioninfo -icon=icon.ico -manifest luzhibo.manifest

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/Baozisoftware/GoldenDaemon"
)

const ver = 2017041400
const p = "录直播"

var port = 12216

var nhta *bool

var htaproc *os.Process

var nt *bool

func main() {
	p := flag.Int("port", port, "WebUI监听端口")
	nopen := flag.Bool("nopenui", false, "不自动打开WebUI")
	nhta = flag.Bool("nhta", false, "禁用hta(仅Windows有效)")
	flag.Bool("d", false, "启用后台运行(仅非Windows有效)")
	nt = flag.Bool("nt", false, "启用无终端交互模式(仅非Windows有效)")
	flag.Parse()
	port = *p
	s := ":" + strconv.Itoa(port)
	if runtime.GOOS != "windows" {
		GoldenDaemon.RegisterTrigger("d", "-nt")
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		time.Sleep(time.Second * 5)
		d, f := filepath.Split(os.Args[0])
		tp := filepath.Join(d, "."+f+".old")
		os.Remove(tp)
	}()
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	if !*nt || runtime.GOOS == "windows" {
		time.Sleep(time.Second * 2)
		go startServer(s)
		if !*nopen {
			openWebUI(!*nhta)
		}
		cmd()
	} else {
		startServer(s)
	}
}

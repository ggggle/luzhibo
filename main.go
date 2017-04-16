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
	"github.com/Baozisoftware/luzhibo/api/getters"
	"github.com/Baozisoftware/luzhibo/workers"
)

const ver = 2017041600
const p = "录直播"

var port = 12216

var nhta *bool

var htaproc *os.Process

var nt *bool

var proxy *string

func main() {
	p := flag.Int("port", port, "WebUI监听端口")
	nopen := flag.Bool("nopenui", false, "不自动打开WebUI")
	nhta = flag.Bool("nhta", false, "禁用hta(仅Windows有效)")
	flag.Bool("d", false, "启用后台运行(仅非Windows有效)")
	nt = flag.Bool("nt", false, "启用无终端交互模式(仅非Windows有效)")
	proxy = flag.String("proxy", "http://127.0.0.1:8880", "代理服务器(如:\"http://127.0.0.1:8888\".)")
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
	getters.Proxy = *proxy
	workers.Proxy = *proxy
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

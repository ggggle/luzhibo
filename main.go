//go:generate goversioninfo -icon=icon.ico -manifest luzhibo.manifest

package main

import (
	"runtime"
	"flag"
	"strconv"
	"fmt"
	"time"
	"os"
	"path/filepath"
	"os/exec"
	"github.com/jpillora/overseer"
)

const ver = 2017040300
const p = "录直播"

var port = 12216

var nhta *bool

var htaproc *os.Process

func main() {
	if runtime.GOOS == "windows" {
		mainFunc()
	} else {
		overseer.Run(overseer.Config{Program: func(state overseer.State) {
			mainFunc()
		}})
	}
}

func mainFunc() {
	p := flag.Int("port", port, "WebUI监听端口")
	nopen := flag.Bool("nopenui", false, "不自动打开WebUI")
	nhta = flag.Bool("nhta", false, "禁用hta(仅Windows有效)")
	d := flag.Bool("d", false, "启用后台运行")
	nt := flag.Bool("nt", false, "启用无终端交互模式")
	flag.Parse()
	port = *p
	s := ":" + strconv.Itoa(port)
	if *d {
		args := os.Args[1:]
		for i, v := range args {
			if v == "-d" || v == "-d=true" {
				args[i] = "-d=false"
			}
		}
		args = append(args, "-nt")
		cmd := exec.Command(os.Args[0], args...)
		cmd.Start()
		if cmd.Process != nil {
			fmt.Printf("[PID]%d\n", cmd.Process.Pid)
			fmt.Printf("[WebUI]\"%s\"\n", s)
		} else {
			fmt.Println("后台启动失败.")
		}
		os.Exit(0)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		time.Sleep(time.Second * 5)
		d, f := filepath.Split(os.Args[0])
		tp := filepath.Join(d, "."+f+".old")
		os.Remove(tp)
	}()
	fmt.Printf("正在\"%s\"处监听WebUI...\n", s)
	if !*nt {
		go startServer(s)
		if !*nopen {
			openWebUI(!*nhta)
		}
		cmd()
	} else {
		startServer(s)
	}
}

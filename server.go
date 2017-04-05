package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"crypto/tls"
	"errors"
	"github.com/Baozisoftware/luzhibo/api"
	"regexp"
	"runtime"
	"encoding/base64"
	"strings"
	nurl "net/url"
	"github.com/inconshreveable/go-update"
	"path/filepath"
)

type checkRet struct {
	Pass    bool
	Has     bool
	Live    bool
	Err     bool
	Path    string
	FileExt string
	Support bool
}

type tasksRet struct {
	Tasks []*taskInfo
	Err   bool
	E     bool
}

type ajaxHandler struct{}

//ServeHTTP 实现接口
func (_ ajaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	act := r.Form.Get("act")
	switch act {
	case "check":
		tr := checkRet{}
		url := r.Form.Get("url")
		oa := api.New(url)
		if oa == nil {
			tr.Pass = false
		} else {
			tr.Pass = true
			i, l, e := oa.GetRoomInfo()
			if e == nil {
				tr.Has = true
				tr.Path = fmt.Sprintf("[%s]%s_%s", oa.Site, i, time.Now().Format("20060102150405"))
				tr.Live = l
				tr.FileExt = oa.FileExt;
				if oa.NeedFFmpeg {
					tr.Support = hasFFmpeg()
				} else {
					tr.Support = true
				}

			} else {
				tr.Err = true
			}
		}
		j, _ := json.Marshal(tr)
		w.Write(j)
		return
	case "add":
		url, m, p, s := r.Form.Get("url"), r.Form.Get("m"), r.Form.Get("path"), r.Form.Get("run")
		url, _ = nurl.QueryUnescape(url)
		p, _ = nurl.QueryUnescape(p)
		mm, ss := m == "true", s == "true"
		if url != "" && p != "" {
			if addTaskEx(url, p, mm, ss) {
				w.Write([]byte("ok"))
				return
			}
		}
	case "addex":
		urls := r.Form.Get("urls")
		urls, _ = nurl.QueryUnescape(urls)
		i := addTasks(urls)
		w.Write([]byte(strconv.Itoa(i)))
	case "del":
		i, d := r.Form.Get("id"), r.Form.Get("f")
		b := d == "true"
		c, e := strconv.Atoi(i)
		if e == nil {
			if delTask(c-1, b) {
				w.Write([]byte("ok"))
				return
			}
		}

	case "start":
		i := r.Form.Get("id")
		if startOrStopTask(i, true) {
			w.Write([]byte("ok"))
			return
		}
	case "stop":
		i := r.Form.Get("id")
		if startOrStopTask(i, false) {
			w.Write([]byte("ok"))
			return
		}
	case "tasks":
		list, o, e := getTaskInfoList()
		r := tasksRet{}
		r.Err = o
		r.Tasks = list
		r.E = e
		j, _ := json.Marshal(r)
		w.Write(j)
		return
	case "exist":
		p := r.Form.Get("path")
		if pp, _ := pathExist(p); pp {
			w.Write([]byte("exist"))
			return
		}
	case "get":
		i, s := r.Form.Get("id"), r.Form.Get("sub")
		ii, e := strconv.Atoi(i)
		if e == nil {
			inf, _ := getTaskInfo(ii - 1)
			fp := inf.Path
			if s != "" {
				fp += "/" + s + "." + inf.FileExt
			}
			pp := inf.Path
			if inf.M {
				if s != "" {
					pp += "_" + s
				}
				pp += "." + inf.FileExt
			}
			w.Header().Add("Content-Disposition", "attachment; filename=\""+nurl.QueryEscape(pp)+"\"")
			w.Header().Add("Content-Type", "video/x-"+inf.FileExt)
			getAct(fp, w)
		}
		return
	case "ver":
		w.Write([]byte(checkUpdate()))
		return
	case "supports":
		lines := api.GetSupports()
		s := strings.Join(lines, "/")
		w.Write([]byte(s))
		return
	case "update":
		b := doUpdate()
		r := "false"
		if b {
			r = "true"
		}
		w.Write([]byte(r))
		return
	}
	w.Write([]byte(""))
}

type uiHandler struct{}

//ServeHTTP 实现接口
func (_ uiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b64 := base64.StdEncoding
	switch r.URL.Path {
	case "/":
		h := ui_main
		r.ParseForm()
		if r.Form.Get("hta") == "true" {
			h = strings.Replace(h, "hta = false", "hta = true", 1)
		}
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte(h))
	case "/hta":
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte(hta))
	case "/favicon.ico":
		w.Header().Add("Content-Type", "image/x-icon")
		data, _ := b64.DecodeString(favicon_ico)
		w.Write(data)
	case "/donate_alipay.png":
		w.Header().Add("Content-Type", "image/png")
		data, _ := b64.DecodeString(donate_alipay_png)
		w.Write(data)
	case "/donate.gif":
		w.Header().Add("Content-Type", "image/gif")
		data, _ := b64.DecodeString(donate_gif)
		w.Write(data)
	case "/donate_wechat.png":
		w.Header().Add("Content-Type", "image/png")
		data, _ := b64.DecodeString(donate_wechat_png)
		w.Write(data)
	case "/bootstrap.min.css":
		w.Header().Add("Content-Type", "text/css")
		data, _ := b64.DecodeString(bootstrap_min_css)
		w.Write(data)
	case "/bootstrap.min.js":
		w.Header().Add("Content-Type", "application/javascript")
		data, _ := b64.DecodeString(bootstrap_min_js)
		w.Write(data)
	case "/jquery.min.js":
		w.Header().Add("Content-Type", "application/javascript")
		data, _ := b64.DecodeString(jquery_min_js)
		w.Write(data)
	case "/flv.min.js":
		w.Header().Add("Content-Type", "application/javascript")
		data, _ := b64.DecodeString(flv_min_js)
		w.Write(data)
	case "/fonts/glyphicons-halflings-regular.woff2":
		w.Header().Add("Content-Type", "application/octet-stream")
		data, _ := b64.DecodeString(glyphicons_halflings_regular_woff2)
		w.Write(data)
	case "/fonts/glyphicons-halflings-regular.eot":
		w.Header().Add("Content-Type", "application/octet-stream")
		data, _ := b64.DecodeString(glyphicons_halflings_regular_eot)
		w.Write(data)
	}
}

func getFile(path string, w http.ResponseWriter) {
	f, e := os.Open(path)
	defer f.Close()
	eof := false
	if e == nil {
		buf := make([]byte, bytes.MinRead)
		for {
			t, e := f.Read(buf)
			if e != nil {
				if e == io.EOF {
					eof = true
				} else {
					break
				}
			}
			_, e = w.Write(buf[:t])
			if e != nil || eof {
				break
			}
		}
	}
}

func getDir(path string, w http.ResponseWriter) {
	files, err := ioutil.ReadDir(path)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() {
				p := path + "/" + f.Name()
				getFile(p, w)
			}
		}
	}
}

func getAct(path string, w http.ResponseWriter) {
	if pe, d := pathExist(path); pe {
		if d {
			getDir(path, w)
		} else {
			getFile(path, w)
		}
	} else {
		w.Write([]byte("no exist"))
	}

}

func startOrStopTask(i string, m bool) bool {
	c, e := strconv.Atoi(i)
	if e != nil {
		return false
	}
	c--
	if m {
		return startTask(c)
	}
	return stopTask(c)
}

func startServer(s string) {
	http.Handle("/", uiHandler{})
	http.Handle("/bootstrap.min.css", uiHandler{})
	http.Handle("/bootstrap.min.js", uiHandler{})
	http.Handle("/jquery.min.css", uiHandler{})
	http.Handle("/flv.min.css", uiHandler{})
	http.Handle("/ajax", ajaxHandler{})
	http.ListenAndServe(s, nil)
	panic("WebUI启动失败.")
}

func httpGet(url string) (data string, err error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	var req *http.Request
	client := http.Client{Transport: tr}
	req, err = http.NewRequest("GET", url, nil)
	if err == nil {
		resp, err := client.Do(req)
		var body []byte
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				data = string(body)
			}
		} else {
			err = errors.New("resp StatusCode is not 200.")
		}
	}
	return
}

func checkUpdate() string {
	r := strconv.Itoa(ver) + "|"
	if updated || updatting {
		r += "null"
	} else {
		data, err := httpGet("https://api.github.com/repos/Baozisoftware/luzhibo/releases/latest")
		if err == nil {
			if data != "" {
				reg, _ := regexp.Compile("Ver (\\d{10})")
				data = reg.FindStringSubmatch(data)[1]
				if v, _ := strconv.Atoi(data); v > ver {
					r += data
				} else {
					r += "null"
				}
			} else {
				r += "null"
			}
		}
	}
	return r
}

var updated = false
var updatting = false

func doUpdate() bool {
	if updatting {
		return false
	}
	if updated {
		return true
	}
	updatting = true
	url := fmt.Sprintf("https://github.com/Baozisoftware/luzhibo/releases/download/latest/luzhibo_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		url += ".exe"
	}
	resp, err := http.Get(url)
	if err != nil {
		updatting = false
		return false
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{OldSavePath: ""})
	if err != nil {
		updatting = false
		return false
	}
	updated = true
	updatting = false
	if runtime.GOOS == "windows" {
		if len(tasks) == 0 {
			if htaproc != nil {
				htaproc.Kill()
			}
			restartSelf()
		}
	}
	return true
}

func restartSelf() {
	args := os.Args
	n := args[0]
	for i, v := range args {
		if v == "-pid" {
			args[i] = ""
			args[i+1] = ""
		}
		if strings.HasPrefix(v, "-pid=") {
			args[i] = ""
		}
	}
	args = append(args, "-pid", strconv.Itoa(os.Getpid()))
	startProc(n, args)
	os.Exit(0)
}

func startProc(name string, args []string) error {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	p, err := filepath.Abs(name)
	if err != nil {
		return err
	}
	_, err = os.StartProcess(p, args, procAttr)
	return err
}

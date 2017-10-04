package workers

import (
    "os/exec"
    "os"
    "github.com/ggggle/luzhibo/api"
    "fmt"
    "bytes"
    "strings"
)

func YoutubeUpload(roomId, fPath string) {
    info, err := os.Stat(fPath)
    if err != nil {
        api.Logger.Print(fPath + " error")
        return
    }
    title := fmt.Sprintf("\"%s-%s\"", roomId, info.ModTime().Format("20060102-1504"))
    cmd := exec.Command("youtube-upload", "--client-secrets", "/root/.client_secret.json",
        "--privacy", "private", "--title", title, fPath)
    w := bytes.NewBuffer(nil)
    cmd.Stderr = w
    cmd.Run()
    uploadRet := string(w.Bytes())
    success := strings.Contains(uploadRet, "Video URL")
    if success {
        api.Logger.Printf("[%s]上传成功", fPath)
        os.Remove(fPath)
    } else {
        api.Logger.Print(uploadRet)
    }
}

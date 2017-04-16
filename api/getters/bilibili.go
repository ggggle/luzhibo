package getters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//bilibili Bilibili直播
type bilibili struct{}

//Site 实现接口
func (i *bilibili) Site() string { return "Bilibili直播" }

//SiteURL 实现接口
func (i *bilibili) SiteURL() string {
	return "http://live.bilibili.com"
}

//SiteIcon 实现接口
func (i *bilibili) SiteIcon() string {
	return i.SiteURL() + "/favicon.ico"
}

//FileExt 实现接口
func (i *bilibili) FileExt() string {
	return "flv"
}

//NeedFFMpeg 实现接口
func (i *bilibili) NeedFFMpeg() bool {
	return false
}

//GetRoomInfo 实现接口
func (i *bilibili) GetRoomInfo(url string) (id string, live bool, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	tmp, err := httpGet(url)
	reg, _ := regexp.Compile("ROOMID = (\\d+)")
	id = reg.FindStringSubmatch(tmp)[1]
	url = "http://live.bilibili.com/live/getInfo?roomid=" + id
	tmp, err = httpGet(url)
	if !strings.Contains(tmp, "\\u623f\\u95f4\\u4e0d\\u5b58\\u5728") {
		live = strings.Contains(tmp, "\"_status\":\"on\"")
	} else {
		err = errors.New("fild get data")
	}
	if id == "" {
		err = errors.New("fail get data")
	}
	return
}

//GetLiveInfo 实现接口
func (i *bilibili) GetLiveInfo(id string) (live LiveInfo, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("fail get data")
		}
	}()
	live = LiveInfo{RoomID: id}
	url := "http://api.live.bilibili.com/live/getInfo?roomid=" + id
	tmp, err := httpGet(url)
	json := *(pruseJSON(tmp).JToken("data"))
	title := json["ROOMTITLE"].(string)
	nick := json["ANCHOR_NICK_NAME"].(string)
	img := json["COVER"].(string)
	cid := json["ROOMID"].(float64)
	url = fmt.Sprintf( "http://live.bilibili.com/api/playurl?player=1&cid=%.0f",cid)
	tmp, err = httpGet(url)
	x, y := strings.Index(tmp, "<url><![CDATA[")+14, strings.LastIndex(tmp, "]]></url>")
	video := tmp[x:y]
	live.LiveNick = nick
	live.RoomTitle = title
	live.LivingIMG = img
	live.RoomDetails = ""
	live.VideoURL = video
	if video == "" {
		err = errors.New("fail get data")
	}
	return
}

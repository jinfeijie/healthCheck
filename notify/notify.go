package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinfeijie/healthCheck/cache"
	"github.com/jinfeijie/healthCheck/echo"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Notify struct{}

var notify *cache.MCache

func NewNotify() *Notify {
	notify = cache.NewMCache()
	return &Notify{}
}

func (n *Notify) Notify(notify []gjson.Result, name, notifyFormat, domain string) {
	mc := cache.NewInstanceMCache()
	var mails []string
	for _, mail := range notify {
		mails = append(mails, mail.String())
	}

	str := fmt.Sprintf(notifyFormat, name, domain, time.Now().Format("2006-1-2 15:04:05"))

	data := &url.URL{}
	u := data.Query()
	u.Add("token", mc.Get("email_token"))

	content := new(struct {
		Title   string
		To      []string
		Content string
	})

	content.To = mails
	content.Title = "网站异常"
	content.Content = str

	jsonStr, err := json.Marshal(content)
	if err != nil {
		echo.Echo(err.Error())
		return
	}

	echo.Echo(mc.Get("email_server") + "?" + u.Encode())
	resp, err := http.Post(mc.Get("email_server")+"?"+u.Encode(), "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		echo.Echo(err.Error())
		return
	}

	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		echo.Echo(err.Error())
		return
	}

	echo.Echo(string(rep))
}

func (n *Notify) PushNow(siteName string) bool {
	if notify.IsExist(siteName) {
		t, err := strconv.Atoi(notify.Get(siteName))
		if err != nil {
			return false
		}

		// 已经过期
		if time.Now().Unix()-int64(t) > 0 {
			return true
		}
		return false
	}
	return true
}

func (n *Notify) Pushed(siteName, expireTime string) {
	notify.Set(siteName, expireTime)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinfeijie/healthCheck/cache"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	mc     *cache.MCache
	notify *cache.MCache
	m      sync.Mutex
)

func init() {
	mc = cache.NewMCache()
	notify = cache.NewMCache()
	mc.LoadData()
	go func() {
		for {
			mc.LoadData()
			<-time.After(time.Second)
		}
	}()
}

func main() {
	for {
		ret := gjson.Parse(mc.Get("site"))
		for _, site := range ret.Array() {
			go func() {
				for i := 0; i < 10; i++ {
					domain := site.Get("domain")
					name := site.Get("name")
					notifyStr := site.Get("notify")
					notifyInterval := site.Get("notify_interval")
					notifyFormat := site.Get("notify_format")
					result := site.Get("result")
					resp, err := http.Get(domain.String())
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					pong, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					resp.Body.Close()
					if string(pong) != result.String() {
						if CanPush(domain.String()) {
							m.Lock()
							if CanPush(domain.String()) {
								Notify(notifyStr.Array(), name.String(), notifyFormat.String(), domain.String())
								notify.Set(domain.String(), strconv.Itoa(int(time.Now().Unix()+notifyInterval.Int())))
							}
							m.Unlock()
						}
					}

					checkInterval, err := strconv.Atoi(mc.Get("check_interval"))

					for i := 0; i < checkInterval; i++ {
						<-time.After(time.Microsecond)
					}
				}
			}()
		}
		<-time.After(time.Second)
	}
}

func CanPush(siteName string) bool {
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

func Notify(notify []gjson.Result, name, notifyFormat, domain string) {
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
		panic(err.Error())
	}

	resp, err := http.Post(mc.Get("email_server")+"?"+u.Encode(), "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err.Error())
	}

	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(rep))

}

package main

import (
	"github.com/jinfeijie/healthCheck/cache"
	"github.com/jinfeijie/healthCheck/echo"
	"github.com/jinfeijie/healthCheck/health"
	message "github.com/jinfeijie/healthCheck/notify"
	"github.com/tidwall/gjson"
	"strconv"
	"sync"
	"time"
)

var (
	mc  *cache.MCache
	m   sync.Mutex
	msg *message.Notify
)

func init() {
	mc = cache.NewInstanceMCache()
	mc.LoadData()

	msg = message.NewNotify()

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
		// 遍历配置监控站点
		for _, site := range ret.Array() {
			go func(site gjson.Result) {
				t := time.Now()
				// 只监控一个时间颗粒
				for time.Now().Sub(t).Seconds() < 60 {
					domain := site.Get("domain")
					name := site.Get("name")
					notifyStr := site.Get("notify")
					notifyInterval := site.Get("notify_interval")
					notifyFormat := site.Get("notify_format")
					result := site.Get("result")

					pong := health.NewHealth(domain.String()).Ping()
					if pong != result.String() {
						if msg.PushNow(domain.String()) {
							m.Lock()
							if msg.PushNow(domain.String()) {
								msg.Notify(notifyStr.Array(), name.String(), notifyFormat.String(), domain.String())
								msg.Pushed(domain.String(), strconv.Itoa(int(time.Now().Unix()+notifyInterval.Int())))
							}
							m.Unlock()
						}
					}

					echo.Echo(domain.String(), pong)

					checkInterval, _ := strconv.Atoi(mc.Get("check_interval"))
					for delay := 0; delay < checkInterval; delay++ {
						<-time.After(time.Millisecond * 100)
					}
				}
			}(site)
		}
		<-time.After(time.Minute)
	}
}

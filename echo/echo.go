package echo

import (
	"fmt"
	"github.com/jinfeijie/healthCheck/cache"
	"time"
)

func Echo(data ...interface{}) {
	c := cache.NewInstanceMCache()
	if c.Get("mod") == "debug" {
		prefix := fmt.Sprintf("[%s] [%s] ", c.Get("app_name"), time.Now().Format("2006-1-2 15:04:05"))
		fmt.Println(prefix, data)
	}
}

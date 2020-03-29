package echo

import (
	"fmt"
	"github.com/go-acme/lego/log"
	"github.com/jinfeijie/healthCheck/cache"
	"runtime"
)

func Echo(data ...interface{}) {
	c := cache.NewInstanceMCache()
	if c.Get("mod") == "debug" {
		var prefix string
		_, file, line, ok := runtime.Caller(1)
		if ok {
			prefix = fmt.Sprintf("[%s] [%s:%d] ", c.Get("app_name"),  file, line)
		} else {
			prefix = fmt.Sprintf("[%s] ", c.Get("app_name"))
		}
		log.Println(prefix, data)
	}
}

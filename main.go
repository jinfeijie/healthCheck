package main

import (
	"fmt"
	"github.com/jinfeijie/healthCheck/cache"
	"time"
)

var (
	mc *cache.MCache
)

func init() {
	mc = cache.NewMCache()
	mc.LoadData()

	go func() {
		for {
			//mc.Dump()
			time.Sleep(time.Second)
		}
	}()
}

func main() {
	fmt.Println(mc.All())
}

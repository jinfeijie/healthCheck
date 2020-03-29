package cache

import (
	"fmt"
	uuid2 "github.com/hashicorp/go-uuid"
	"testing"
	"time"
)

func TestNewInstanceMCache(t *testing.T) {
	mc := NewInstanceMCache()
	mc.LoadData()
	mc2 := NewInstanceMCache()

	if mc.Get("app_name") != mc2.Get("app_name") {
		t.Logf("mc:%s  mc2:%s", mc.Get("app_name"), mc2.Get("app_name"))
		t.Fail()
	}

	t.Logf("mc:%s  mc2:%s", mc.Get("app_name"), mc2.Get("app_name"))
}

func TestNewMCache(t *testing.T) {
	mc := NewMCache()
	mc.LoadData()
	mc2 := NewMCache()

	if mc.Get("app_name") != mc2.Get("app_name") {
		t.Logf("mc:%s  mc2:%s", mc.Get("app_name"), mc2.Get("app_name"))
		t.Fail()
	}

	t.Logf("mc:%s  mc2:%s", mc.Get("app_name"), mc2.Get("app_name"))
}

func TestNewMCache2(t *testing.T) {
	mc := NewMCache()
	mc.LoadData()
	mc2 := NewMCache()
	mc2.LoadData()

	if mc.Get("app_name") != mc2.Get("app_name") {
		t.Logf("mc:%s  mc2:%s", mc.Get("app_name"), mc2.Get("app_name"))
		t.Fail()
	}

	t.Logf("mc:%s  mc2:%s", mc.Get("app_name"), mc2.Get("app_name"))
}

func TestInstanceCacheAllMethod(t *testing.T) {
	mc := NewInstanceMCache()
	mc.Set("test", "test")
	t.Log(mc.Get("test"))

	mc.Set("test1", "test1")
	t.Log(mc.Get("test1"))

	t.Log(mc.IsExist("test"))

	t.Log(mc.All())

	mc.Del("test")

	t.Log(mc.IsExist("test"))

	t.Log(mc.All())
}

func TestCacheAllMethod(t *testing.T) {
	mc := NewMCache()
	mc.Set("test", "test")
	t.Log(mc.Get("test"))

	mc.Set("test1", "test1")
	t.Log(mc.Get("test1"))

	t.Log(mc.IsExist("test"))

	t.Log(mc.All())

	mc.Del("test")

	t.Log(mc.IsExist("test"))

	t.Log(mc.All())

}

func TestBatchAction(t *testing.T) {
	c := NewInstanceMCache()
	c.LoadData()

	go func(c *MCache) {
		for {
			fmt.Println("Get", c.Get("test"))
		}
	}(c)

	go func(c *MCache) {
		for {
			uuid, err := uuid2.GenerateUUID()
			if err != nil {
				fmt.Println(err.Error())
				uuid = ""
			}
			c.Set("test", uuid)
			fmt.Println("Set")
		}
	}(c)

	go func(c *MCache) {
		for {
			c.IsExist("test")
			fmt.Println("IsExist")
		}
	}(c)

	go func(c *MCache) {
		for {
			c.Del("test")
			fmt.Println("Del")
		}
	}(c)

	go func(c *MCache) {
		for {
			for _, _ = range c.All() {

			}
			fmt.Println("Range")
		}
	}(c)


	go func(c *MCache) {
		for {
			c.Len()
			fmt.Println("Len")
		}
	}(c)

	go func(c *MCache) {
		for {
			c.LoadData()
			fmt.Println("LoadData")
		}
	}(c)

	<- time.After(time.Minute * 5)
}

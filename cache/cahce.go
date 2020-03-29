package cache

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"sync"
)

type MCache struct {
	Data map[string]string
	sync.RWMutex
}

type Interface interface {
	Set(key, value string) bool
	Get(key string) string
	IsExist(key string) bool
	Del(key string) bool
	All() map[string]string
	Len() int
	Dump()
	LoadData()
}

var once sync.Once
var mc *MCache

func NewMCache() *MCache {
	data := make(map[string]string)
	return &MCache{
		Data: data,
	}
}

func NewInstanceMCache() *MCache {
	data := make(map[string]string)
	once.Do(func() {
		mc = &MCache{
			Data: data,
		}
	})
	return mc
}

func (mc *MCache) Set(key, value string) bool {
	mc.Lock()
	defer mc.Unlock()

	mc.Data[key] = value
	return true
}

func (mc *MCache) Get(key string) string {
	mc.RLock()
	defer mc.RUnlock()

	if val, ok := mc.Data[key]; ok {
		return val
	}
	return ""
}

func (mc *MCache) IsExist(key string) bool {
	mc.RLock()
	defer mc.RUnlock()

	if _, ok := mc.Data[key]; ok {
		return true
	}
	return false
}

func (mc *MCache) Del(key string) bool {
	mc.Lock()
	defer mc.Unlock()

	if _, ok := mc.Data[key]; ok {
		delete(mc.Data, key)
		return true
	}
	return false
}

func (mc *MCache) All() map[string]string {
	mc.RLock()
	defer mc.RUnlock()

	allData := make(map[string]string)
	for key, value := range mc.Data {
		allData[key] = value
	}
	return allData
}
func (mc *MCache) Len() int {
	mc.RLock()
	defer mc.RUnlock()
	return len(mc.Data)
}

func (mc *MCache) Dump() {
	dumpData := mc.All()

	data, err := json.Marshal(dumpData)
	if err != nil {
		panic(err.Error())
	}
	f, _ := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	f.Truncate(0)
	_, _ = f.Write(data)
}

func (mc *MCache) LoadData() {
	_, err := os.Stat("config.json")
	if err != nil {
		panic(err.Error())
	}

	f, _ := os.OpenFile("config.json", os.O_RDWR, 0666)
	content, err := ioutil.ReadAll(f)
	str := gjson.Parse(string(content))
	str.ForEach(func(key, value gjson.Result) bool {
		mc.Data[key.String()] = value.String()
		return true
	})
}

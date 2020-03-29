package health

import (
	"github.com/jinfeijie/healthCheck/echo"
	"io/ioutil"
	"net/http"
)

type Health struct {
	Domain string
}

//NewHealth
func NewHealth(domain string) *Health {
	return &Health{Domain: domain}
}

func (h *Health) Ping() string {
	resp, err := http.Get(h.Domain)
	defer resp.Body.Close()
	if err != nil {
		echo.Echo(err.Error())
		return ""
	}

	pong, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		echo.Echo(err.Error())
		return ""
	}
	return string(pong)
}

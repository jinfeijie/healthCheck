package echo

import (
	"github.com/jinfeijie/healthCheck/cache"
	"testing"
)

func init() {
	// load cache
	cache.NewInstanceMCache().LoadData()
}

func TestEcho(t *testing.T) {
	Echo(1,2,3,4,5)
}

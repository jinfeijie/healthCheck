package health

import "testing"

func TestHealth_Ping(t *testing.T) {
	health := NewHealth("https://jinfeijie.cn/a.php")
	t.Log(health.Ping())
}

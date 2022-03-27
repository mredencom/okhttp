package okhttp

import (
	"github.com/mredencom/okhttp/log"
	"testing"
)

func TestGet(t *testing.T) {
	okhttp := log.New()
	t.Log(okhttp.Get("http://httpbin.org"))
}

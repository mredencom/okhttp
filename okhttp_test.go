package okhttp

import "testing"

func TestGet(t *testing.T) {
	okhttp := New()
	t.Log(okhttp.Get("http://httpbin.org"))
}

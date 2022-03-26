package okhttp

import "net/http"

// doHeader do http request header
func doHeader(headers http.Header) H {
	var h = H{}
	if headers == nil {
		return h
	}
	for k := range headers {
		h[k] = headers.Get(k)
	}
	return h
}

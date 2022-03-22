package okhttp

import "net/http"

// addHeaders 增加请求头
func addHeaders(r *http.Request, headers http.Header) {
	for key, header := range headers {
		for _, value := range header {
			r.Header.Add(key, value)
		}
	}
}

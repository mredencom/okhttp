package okhttp

import (
	"encoding/json"
	"net/http"
)

// Response r
type Response struct {
	request *Request
	headers http.Header
	cookies []*http.Cookie
	status  int
	body    []byte
}

// GetCookies returns response cookies slice
func (r *Response) GetCookies() []*http.Cookie {
	return r.cookies
}

// GetBody returns response body
func (r *Response) GetBody() []byte {
	return r.body
}

// String returns response body as string
func (r *Response) String() string {
	return string(r.GetBody())
}

// GetJSON unmarshal JSON response to struct
func (r *Response) GetJSON(v interface{}) error {
	return json.Unmarshal(r.body, v)
}

// GetHeaders return response headers
func (r *Response) GetHeaders() H {
	return doHeader(r.headers)
}

// GetRequest returns initial request
func (r *Response) GetRequest() *Request {
	return r.request
}

// GetStatus returns response status code
func (r *Response) GetStatus() int {
	return r.status
}

// GetHeader returns response header by name
func (r *Response) GetHeader(key string) string {
	return r.headers.Get(key)
}

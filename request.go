package okhttp

import "net/http"

type RequestOption func(req *RequestCtx)

type RequestCtx struct {
	r     http.Request
	Debug bool

	http.Client
}

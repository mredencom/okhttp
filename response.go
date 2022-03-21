package okhttp

import "net/http"

type ResponseOption func(req *ResponseCtx)

type ResponseCtx struct {
	r http.Response
}

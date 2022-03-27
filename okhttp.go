package okhttp

import (
	"fmt"
	"github.com/mredencom/okhttp/log"
	"net/http"
	"net/url"
	"strings"
)

type Handler func(*Response)

// NewRequest 建立一个请求
func NewRequest(method, uri string) (r *Request, err error) {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		method = http.MethodGet
		break
	case http.MethodPost:
		method = http.MethodPost
		break
	case http.MethodPut:
		method = http.MethodPut
		break
	case http.MethodDelete:
		method = http.MethodDelete
		break
	case http.MethodHead:
		method = http.MethodHead
		break
	case http.MethodPatch:
		method = http.MethodPatch
		break
	case http.MethodOptions:
		method = http.MethodOptions
		break
	case http.MethodTrace:
		method = http.MethodTrace
		break
	case http.MethodConnect:
		method = http.MethodConnect
		break
	default:
		err = NoMatchHttpMethod
		break
	}
	// parse url
	parse, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	// header
	header := http.Header{}
	header.Set("User-Agent", fmt.Sprintf("Okhttp/%s", Version))
	r = &Request{
		method:        method,
		url:           parse,
		header:        header,
		allowRedirect: true,
		debug:         false,
		isPrintBody:   false,
		l:             log.NewLogger(0),
	}
	return
}

// Get created a get request
func Get(uri string) (*Request, error) {
	return NewRequest(http.MethodGet, uri)
}

// Post created a post request
func Post(uri string) (*Request, error) {
	return NewRequest(http.MethodPost, uri)
}

// Put created a put request
func Put(uri string) (*Request, error) {
	return NewRequest(http.MethodPut, uri)
}

// Delete created a delete request
func Delete(uri string) (*Request, error) {
	return NewRequest(http.MethodDelete, uri)
}

// Head created a head request
func Head(uri string) (*Request, error) {
	return NewRequest(http.MethodHead, uri)
}

// Patch created a patch request
func Patch(uri string) (*Request, error) {
	return NewRequest(http.MethodPatch, uri)
}

// Options created a options request
func Options(uri string) (*Request, error) {
	return NewRequest(http.MethodOptions, uri)
}

// Trace created a trace request
func Trace(uri string) (*Request, error) {
	return NewRequest(http.MethodTrace, uri)
}

// Connect created a connect request
func Connect(uri string) (*Request, error) {
	return NewRequest(http.MethodConnect, uri)
}

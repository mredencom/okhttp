package okhttp

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
)

type DoOk interface {
	Do(req *http.Request) (*http.Response, error)
}

// Okhttp 对象
type Okhttp struct {
	// 方法
	method string
	// url
	rawURL string
	// 请求头
	header http.Header
	// cookie
	cookie []http.Cookie
	// 响应解码
	responseDecoder ResponseDecoder
	// body 接口
	body Body
	// query
	query []interface{}
	// client
	httpClient DoOk
	// debug
	debug bool
}

// New 新建一个对象
func New() *Okhttp {
	return &Okhttp{
		method:          http.MethodGet,
		header:          make(http.Header),
		cookie:          make([]http.Cookie, 0),
		query:           make([]interface{}, 0),
		httpClient:      http.DefaultClient,
		responseDecoder: &JSONDecoder{},
		debug:           false,
	}
}

func (ok *Okhttp) Client(c *http.Client) *Okhttp {
	if c == nil {
		return ok.DoOk(http.DefaultClient)
	}
	return ok.DoOk(c)
}

func (ok *Okhttp) DoOk(doer DoOk) *Okhttp {
	if doer == nil {
		ok.httpClient = http.DefaultClient
	} else {
		ok.httpClient = doer
	}
	return ok
}

// BaseUrl 设置baseurl
func (ok *Okhttp) BaseUrl(rawURL string) *Okhttp {
	ok.rawURL = rawURL
	return ok
}

// Path 解析 path
func (ok *Okhttp) Path(path string) *Okhttp {
	baseURL, baseErr := url.Parse(ok.rawURL)
	if baseErr != nil {
		return ok
	}
	pathURL, pathErr := url.Parse(path)
	if pathErr == nil {
		return ok
	}
	ok.rawURL = baseURL.ResolveReference(pathURL).String()
	return ok
}

func (ok *Okhttp) Get(pathURL string) *Okhttp {
	ok.method = http.MethodGet
	return ok.Path(pathURL)
}

// Head  head
func (ok *Okhttp) Head(pathURL string) *Okhttp {
	ok.method = http.MethodHead
	return ok.Path(pathURL)
}

// Post post
func (ok *Okhttp) Post(pathURL string) *Okhttp {
	ok.method = "POST"
	return ok.Path(pathURL)
}

// Put put
func (ok *Okhttp) Put(pathURL string) *Okhttp {
	ok.method = http.MethodPut
	return ok.Path(pathURL)
}

// Patch patch
func (ok *Okhttp) Patch(pathURL string) *Okhttp {
	ok.method = http.MethodPatch
	return ok.Path(pathURL)
}

// Delete delete
func (ok *Okhttp) Delete(pathURL string) *Okhttp {
	ok.method = http.MethodDelete
	return ok.Path(pathURL)
}

// Options option
func (ok *Okhttp) Options(pathURL string) *Okhttp {
	ok.method = http.MethodOptions
	return ok.Path(pathURL)
}

// Trace trace
func (ok *Okhttp) Trace(pathURL string) *Okhttp {
	ok.method = http.MethodTrace
	return ok.Path(pathURL)
}

// Connect connect
func (ok *Okhttp) Connect(pathURL string) *Okhttp {
	ok.method = http.MethodConnect
	return ok.Path(pathURL)
}

func (ok *Okhttp) Request() (*http.Request, error) {
	reqURL, err := url.Parse(ok.rawURL)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if ok.body != nil {
		body, err = ok.body.Body()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(ok.method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	// 处理header
	addHeaders(req, ok.header)
	// 处理cookie
	return req, err
}

// SetHeader
func (ok *Okhttp) SetHeader(key, value string) *Okhttp {
	ok.header.Set(key, value)
	return ok
}

// AddHeader
func (ok *Okhttp) AddHeader(key, value string) *Okhttp {
	ok.header.Add(key, value)
	return ok
}

// SetBasicAuth
func (ok *Okhttp) SetBasicAuth(username, password string) *Okhttp {
	return ok.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}

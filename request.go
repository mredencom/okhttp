package okhttp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/mredencom/okhttp/log"
	"golang.org/x/net/publicsuffix"
)

// DefaultRequestTimeOut set a default request time out
var DefaultRequestTimeOut = 5 * time.Second

type Request struct {
	method        string
	url           *url.URL
	header        http.Header
	cookies       []*http.Cookie
	body          io.Reader
	timeout       time.Duration
	proxy         func(*http.Request) (*url.URL, error)
	allowRedirect bool
	debug         bool
	isPrintBody   bool
	l             *log.Logger
}

//SetDebug set debug mode
func (r *Request) SetDebug(d bool) *Request {
	r.debug = d
	return r
}

// SetPrintBody set debug mode
// rely Debug
func (r *Request) SetPrintBody(d bool) *Request {
	r.isPrintBody = d
	return r
}

// Method get http method
func (r *Request) Method() string {
	return r.method
}

// URL url object
func (r *Request) URL() *url.URL {
	return r.url
}

// URLString url string
func (r *Request) URLString() string {
	return r.url.String()
}

// SetTimeOut set default request time
func (r *Request) SetTimeOut(d time.Duration) *Request {
	if d > 0 {
		r.timeout = d
		return r
	}
	r.timeout = DefaultRequestTimeOut
	return r
}

// SetRedirects set default request allow redirects
func (r *Request) SetRedirects(i bool) *Request {
	r.allowRedirect = i
	return r
}

// GetHeaders get all request header
func (r *Request) GetHeaders() H {
	return doHeader(r.header)
}

// SetHeader set request header
func (r *Request) SetHeader(key, value string) *Request {
	r.header.Set(key, value)
	return r
}

// SetHeaders multi set request header
func (r *Request) SetHeaders(headers map[string]string) *Request {
	if headers == nil {
		return r
	}
	for k, v := range headers {
		r.SetHeader(k, v)
	}
	return r
}

// AddHeader add request header
func (r *Request) AddHeader(key, value string) *Request {
	r.header.Add(key, value)
	return r
}

// SetUserAgent set user-agent
func (r *Request) SetUserAgent(value string) *Request {
	r.SetHeader("User-Agent", value)
	return r
}

// SetReferer set a referer request header
func (r *Request) SetReferer(referer string) *Request {
	r.SetHeader("Referer", referer)
	return r
}

// SetProxy set a referer request header
func (r *Request) SetProxy(proxyURL string) *Request {
	parse, err := url.Parse(proxyURL)
	if err != nil {
		panic("illegal url")
	}
	r.proxy = http.ProxyURL(parse)
	return r
}

// SetCookie set a cookie request header
func (r *Request) SetCookie(cookie *http.Cookie) *Request {
	r.cookies = append(r.cookies, cookie)
	return r
}

// SetBody sets request body
func (r *Request) SetBody(body io.Reader) *Request {
	r.body = body
	return r
}

// SetForm sets request form and returns response
func (r *Request) SetForm(v url.Values) *Request {
	r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	// todo err
	return r.SetBody(bytes.NewBuffer([]byte(v.Encode())))
}

// SetJSON sets request JSON and returns response
func (r *Request) SetJSON(v interface{}) *Request {

	r.SetHeader("Content-Type", "application/json")
	body, err := json.Marshal(v)
	if err != nil {
		panic("json encode err " + err.Error())
	}
	return r.SetBody(bytes.NewBuffer(body))
}

// SetBasicAuth set username and password to header
func (r *Request) SetBasicAuth(username, password string) *Request {

	r.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
	return r
}

// Do returns response
func (r *Request) Do() (*Response, error) {

	client, err := r.client()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(r.method, r.url.String(), r.body)
	if err != nil {
		return nil, err
	}
	for k := range r.header {
		request.Header.Add(k, r.header.Get(k))
	}

	if r.debug {
		dumpRequest, _ := httputil.DumpRequest(request, r.isPrintBody)
		r.l.Info(string(dumpRequest))
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return r.response(response)
}

// client create a request client
func (r *Request) client() (*http.Client, error) {

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	if err != nil {
		return nil, err
	}
	jar.SetCookies(r.url, r.cookies)

	client := &http.Client{
		Transport: &http.Transport{
			// 设置代理
			Proxy: r.proxy,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Jar:     jar,
		Timeout: r.timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !r.allowRedirect {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	if len(r.cookies) > 0 {
		client.Jar.SetCookies(r.url, r.cookies)
	}

	return client, nil
}

// response response data
func (r *Request) response(response *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := &Response{
		request: r,
		status:  response.StatusCode,
		headers: response.Header,
		cookies: response.Cookies(),
		body:    body,
	}

	if r.debug {
		dumpResponse, _ := httputil.DumpResponse(response, r.isPrintBody)
		r.l.Info(string(dumpResponse))
	}

	return res, nil
}

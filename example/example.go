package main

import (
	"github.com/mredencom/okhttp"
	"github.com/mredencom/okhttp/log"
)

func main() {
	//RequestGet("http://httpbin.org/get")
	RequestRedirects("http://httpbin.org/absolute-redirect/2")
}

func RequestGet(url string) {
	getR, _ := okhttp.Get(url)
	getR.SetDebug(true)
	resp, _ := getR.Do()
	var s = GetResponse{}
	resp.GetJSON(&s)

	log.Println(s)
}

type GetResponse struct {
	Args struct {
	} `json:"args"`
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		Host           string `json:"Host"`
		Referer        string `json:"Referer"`
		UserAgent      string `json:"User-Agent"`
		XAmznTraceID   string `json:"X-Amzn-Trace-Id"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

//RequestRedirects http://httpbin.org/absolute-redirect/1
func RequestRedirects(url string) {
	getR, _ := okhttp.Get(url)
	//getR.SetDebug(true)
	getR.SetHeader("content-type", "text/plain; charset=utf-8")
	getR.SetRedirects(false)
	resp, _ := getR.Do()
	//log.Println(resp.String())
	log.Println(resp.GetStatus())
}

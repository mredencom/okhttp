package okhttp

import (
	"net/http"
	"time"
)

type ClientOpt func(client *Client)

type Client struct {
	BaseUri string
	TimeOut time.Duration
	Cookies http.CookieJar
	Proxy   string
}

var DefaultClient = &Client{}

func NewClient(opt ...ClientOpt) *Client {
	if opt == nil {
		return DefaultClient
	}
	return &Client{}
}

func TimeOut(d time.Duration) ClientOpt {
	return func(client *Client) {
		client.TimeOut = d
	}
}

// Cookies 建立cookie
func Cookies(cookies Ok) ClientOpt {
	return func(client *Client) {

	}
}

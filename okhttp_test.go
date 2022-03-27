package okhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Headers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf(`Method should be %s, %s given`, http.MethodDelete, r.Method)
		}

		customHeader := r.Header.Get("X-Request-Header")
		userAgent := r.Header.Get("User-Agent")
		referer := r.Header.Get("Referer")

		if customHeader != "Value" {
			t.Errorf(`Custom header "X-Request-Header" should be "%s", "%s" given`, "Value", customHeader)
		}

		if userAgent != "Test" {
			t.Errorf(`User agent should be "%s", "%s" given`, "Test", userAgent)
		}

		if referer != "http://foo.bar/fizz?buz=baz" {
			t.Errorf(`Referer should be "%s", "%s" given`, "http://foo.bar/fizz?buz=baz", referer)
		}

		w.Header().Set("X-Response-Header", "Bite me")
	}))
	defer ts.Close()

	res, err := Delete(ts.URL)
	resp, _ := res.SetHeader("X-Request-Header", "Value").
		SetUserAgent("Test").
		SetReferer("http://foo.bar/fizz?buz=baz").
		Do()

	if err != nil {
		t.Error(err)
	}

	responseHeader := resp.GetHeader("X-Response-Header")
	if responseHeader != "Bite me" {
		t.Errorf(`Response header should be "%s", "%s" given`, "Bite me", responseHeader)
	}
}

func Test_Cookies(t *testing.T) {
	requestCookie := &http.Cookie{
		Name:  "RequestCookie",
		Value: "Some value",
	}

	responseCookie := &http.Cookie{
		Name:  "ResponseCookie",
		Value: "Another value",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf(`Method should be %s, %s given`, http.MethodGet, r.Method)
		}

		cookies := r.Cookies()
		if len(cookies) != 1 {
			t.Errorf("Request should contain 1 cookie, %d given", len(cookies))
		}

		if cookies[0].Name != requestCookie.Name {
			t.Errorf(`Request cookie name should be "%s", "%s" given`, requestCookie.Name, cookies[0].Name)
		}

		if cookies[0].Value != requestCookie.Value {
			t.Errorf(`Request cookie value should be "%s", "%s" given`, requestCookie.Value, cookies[0].Value)
		}

		http.SetCookie(w, responseCookie)
	}))
	defer ts.Close()

	res, err := Get(ts.URL)
	resp, _ := res.SetCookie(requestCookie).Do()
	if err != nil {
		t.Error(err)
	}

	cookies := resp.GetCookies()
	if len(cookies) != 1 {
		t.Errorf("Response should contain 1 cookie, %d given", len(cookies))
	}

	if cookies[0].Name != responseCookie.Name {
		t.Errorf(`Response cookie name should be "%s", "%s" given`, responseCookie.Name, cookies[0].Name)
	}
}

func Test_Body(t *testing.T) {
	requestBody := "Ping"
	responseBody := "Pong"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf(`Method should be %s, %s given`, http.MethodPut, r.Method)
		}

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(reqBody) != requestBody {
			t.Errorf(`Request body should be "%s", "%s" given`, requestBody, string(reqBody))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
	defer ts.Close()

	body := &bytes.Buffer{}
	body.Write([]byte(requestBody))
	res, err := Put(ts.URL)
	if err != nil {
		t.Error(err)
	}
	resp, err := res.SetBody(body).Do()
	if err != nil {
		t.Error(err)
	}

	if resp.String() != responseBody {
		t.Errorf(`Response body should be "%s", "%s" given`, responseBody, resp.String())
	}

	if resp.GetStatus() != http.StatusOK {
		t.Errorf(`Response status should be "%d", "%d" given`, http.StatusOK, resp.GetStatus())
	}
}

type testStruct struct {
	Foo  string `json:"foo"`
	Fizz int    `json:"bar"`
}

func Test_JSON(t *testing.T) {
	testStructValue := testStruct{
		Foo:  "bar",
		Fizz: 42,
	}

	testJSON, err := json.Marshal(&testStructValue)
	if err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf(`Method should be %s, %s given`, http.MethodPost, r.Method)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf(`Content type should be "%s", "%s" given`, "application/json", contentType)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if string(body) != string(testJSON) {
			t.Errorf(`Body should be "%s", "%s" given`, string(testJSON), string(body))
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(testJSON)
	}))
	defer ts.Close()

	res, err := Post(ts.URL)
	resp, err := res.SetJSON(&testStructValue).Do()
	if err != nil {
		t.Error(err)
	}

	if resp.GetStatus() != http.StatusBadRequest {
		t.Errorf(`Response status should be "%d", "%d" given`, http.StatusBadRequest, resp.GetStatus())
	}

	var resValue testStruct
	err = resp.GetJSON(&resValue)
	if err != nil {
		t.Error(err)
	}

	if resValue.Foo != testStructValue.Foo {
		t.Errorf(`Should be "%s"", "%s" given`, testStructValue.Foo, resValue.Foo)
	}

	if resValue.Fizz != testStructValue.Fizz {
		t.Errorf(`Should be %d, %d given`, testStructValue.Fizz, resValue.Fizz)
	}
}

func Test_Form(t *testing.T) {
	form := url.Values{}
	form.Add("Foo", "Bar")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf(`Method should be %s, %s given`, http.MethodPost, r.Method)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			t.Errorf(`Content type should be "%s", "%s" given`, "application/x-www-form-urlencoded", contentType)
		}

		val := r.FormValue("Foo")
		if val != "Bar" {
			t.Errorf(`Form value of "Foo" should be "%s", "%s" given`, "Bar", val)
		}

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	res, err := Post(ts.URL)
	resp, err := res.SetForm(form).Do()
	if err != nil {
		t.Error(err)
	}

	if resp.GetStatus() != http.StatusBadRequest {
		t.Errorf(`Response status should be "%d", "%d" given`, http.StatusBadRequest, resp.GetStatus())
	}
}

func Test_Debug(t *testing.T) {
	userAgent := "okhttp"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	body := &bytes.Buffer{}
	body.Write([]byte("okhttp"))
	res, err := Delete(ts.URL)
	resp, err := res.SetUserAgent(userAgent).SetDebug(true).SetBody(body).Do()

	if err != nil {
		t.Error(err)
	}

	if resp.GetStatus() != http.StatusOK {
		t.Errorf(`Response status should be "%d", "%d" given`, http.StatusBadRequest, resp.GetStatus())
	}
}

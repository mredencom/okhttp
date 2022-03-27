package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/mredencom/okhttp"
	"github.com/mredencom/okhttp/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"strconv"
	"time"
)

func main() {
	//RequestGet("http://httpbin.org/get")
	//RequestRedirects("http://httpbin.org/absolute-redirect/2")
	uploadFile()
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

func uploadFile() {
	flag.Parse()

	positionalArgs := flag.Args()
	if len(positionalArgs) == 0 {
		log.Fatal("This program requires at least 1 positional argument.")
	}
	// Metadata content.
	metadata := `{"title": "hello world", "description": "Multipart related upload test"}`

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Metadata part.
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json")
	metadataHeader.Set("Content-ID", "metadata")
	part, err := writer.CreatePart(metadataHeader)
	if err != nil {
		log.Fatal("Error writing metadata headers: %v", err)
	}
	part.Write([]byte(metadata))

	// Media Files.
	for _, mediaFilename := range positionalArgs {
		mediaData, errRead := ioutil.ReadFile(mediaFilename)
		if errRead != nil {
			log.Fatal("Error reading media file: %v", errRead)
		}
		mediaHeader := textproto.MIMEHeader{}
		mediaHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\".", mediaFilename))
		mediaHeader.Set("Content-ID", "media")
		mediaHeader.Set("Content-Filename", mediaFilename)

		mediaPart, err := writer.CreatePart(mediaHeader)
		if err != nil {
			log.Fatal("Error writing media headers: %v", errRead)
		}

		if _, err := io.Copy(mediaPart, bytes.NewReader(mediaData)); err != nil {
			log.Fatal("Error writing media: %v", errRead)
		}
	}

	// Close multipart writer.
	if err := writer.Close(); err != nil {
		log.Fatal("Error closing multipart writer: %v", err)
	}

	// Request Content-Type with boundary parameter.
	contentType := fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary())
	post, _ := okhttp.Post("http://localhost:8080/upload")
	post.SetHeader("Content-type", contentType).
		SetHeader("Accept", "*/*").
		SetHeader("Content-type", strconv.Itoa(body.Len())).
		SetDebug(true).
		SetTimeOut(100 * time.Second)
	do, _ := post.SetBody(body).Do()

	log.Println(do.String())
}

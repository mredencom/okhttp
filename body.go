package okhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	contentType     = "Content-Type"
	textContentType = "application/text"
	jsonContentType = "application/json"
	formContentType = "application/x-www-form-urlencoded"
)

type Body interface {
	Body() (io.Reader, error)
	ContentType() string
}

// bodyProvider provides the wrapped body value as a Body for reqests.
type textBody struct {
	body io.Reader
}

func (p textBody) ContentType() string {
	return textContentType
}

func (p textBody) Body() (io.Reader, error) {
	return p.body, nil
}

// json impl
type jsonBody struct {
	payload interface{}
}

func (p jsonBody) ContentType() string {
	return jsonContentType
}

func (p jsonBody) Body() (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(p.payload)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

//
type formBody struct {
	payload interface{}
}

func (p formBody) ContentType() string {
	return formContentType
}
func (p formBody) Body() (io.Reader, error) {
	values, err := query.Values(p.payload)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(values.Encode()), nil
}

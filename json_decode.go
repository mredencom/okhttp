package okhttp

import (
	"encoding/json"
	"io"
)

type JSONDecoder struct {
	value interface{}
}

func NewJsonDecoder() *JSONDecoder {
	return &JSONDecoder{}
}

func (d *JSONDecoder) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(d.value)
}

func (d *JSONDecoder) Value() interface{} {
	return d.value
}

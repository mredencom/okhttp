package okhttp

import (
	"encoding/xml"
	"io"
)

type XMLDecoder struct {
	value interface{}
}

func (d *XMLDecoder) Decode(reader io.Reader) error {
	return xml.NewDecoder(reader).Decode(d.value)
}

func (d *XMLDecoder) Value() interface{} {
	return d.value
}

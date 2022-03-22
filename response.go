package okhttp

import "io"

// ResponseDecoder 响应解码接口
type ResponseDecoder interface {
	Decode(r io.Reader) error
	Value() interface{}
}

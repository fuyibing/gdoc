// author: wsfuyibing <websearch@163.com>
// date: 2020-12-13

package tpl

import (
	"github.com/fuyibing/gdoc/base"
)

type (
	// Response
	// 请求结果.
	Response struct {
		Name     string    `json:"name"`
		Request  *Request  `json:"originalRequest"`
		Language string    `json:"_postman_previewlanguage"`
		Header   []*Header `json:"header,omitempty"`
		Cookie   []*Header `json:"cookie,omitempty"`
		Body     string    `json:"body"`
	}
)

func NewResponse() *Response {
	return (&Response{
		Language: "json",
	}).init()
}

func (o *Response) End(key string) {
	o.Body = string(base.Mapper.LoadTmpCode(key))
}

func (o *Response) init() *Response {
	return o
}

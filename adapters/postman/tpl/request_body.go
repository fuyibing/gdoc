// author: wsfuyibing <websearch@163.com>
// date: 2020-12-15

package tpl

import (
	"github.com/fuyibing/gdoc/base"
)

type (
	RequestBody struct {
		Mode RequestMode `json:"mode"`

		UrlEncoded []map[string]interface{} `json:"urlencoded,omitempty"`

		Raw     string                 `json:"raw,omitempty"`
		Options map[string]interface{} `json:"options,omitempty"`
	}

	RequestMode string
)

func NewRequestBody() *RequestBody {
	return (&RequestBody{}).init()
}

func (o *RequestBody) End(key string) {
	if o.Mode == RequestModeRaw {
		o.Raw = string(base.Mapper.LoadTmpCode(key))

		o.Options = map[string]interface{}{
			"raw": map[string]string{
				"language": "json",
			},
		}
	}
}

func (o *RequestBody) init() *RequestBody {
	return o
}

const (
	RequestModeRaw        RequestMode = "raw"
	RequestModeUrlEncoded RequestMode = "urlencoded"
)

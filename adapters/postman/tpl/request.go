// author: wsfuyibing <websearch@163.com>
// date: 2022-12-13

package tpl

type (
	// Request
	// 请求参数.
	Request struct {
		Description string `json:"description"`
		Method      string `json:"method"`

		Body   *RequestBody `json:"body,omitempty"`
		Header []*Header    `json:"header,omitempty"`
		Url    *RequestUrl  `json:"url,omitempty"`
	}
)

func NewRequest() *Request {
	return &Request{
		Header: make([]*Header, 0),
	}
}

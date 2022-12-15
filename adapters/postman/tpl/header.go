// author: wsfuyibing <websearch@163.com>
// date: 2022-12-13

package tpl

type (
	Header struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	}
)

func NewHeader() *Header {
	return &Header{
		Type: "text",
	}
}

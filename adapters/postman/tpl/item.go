// author: wsfuyibing <websearch@163.com>
// date: 2020-12-13

package tpl

type (
	// Item
	// 接口条目.
	Item struct {
		Description string `json:"description,omitempty"`
		Name        string `json:"name,omitempty"`

		Request  *Request    `json:"request,omitempty"`
		Response []*Response `json:"response,omitempty"`

		Item []*Item `json:"item,omitempty"`
	}
)

func NewItem() *Item {
	return &Item{}
}

func (o *Item) Add(items ...*Item) {
	if o.Item == nil {
		o.Item = make([]*Item, 0)
	}
	o.Item = append(o.Item, items...)
}

func (o *Item) SetResponse(response *Response) {
	if o.Response == nil {
		o.Response = make([]*Response, 0)
	}
	o.Response = append(o.Response, response)
}

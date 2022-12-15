// author: wsfuyibing <websearch@163.com>
// date: 2022-12-13

package tpl

type (
	// Postman
	// 模板结构体.
	Postman struct {
		Info *Info   `json:"info"`
		Item []*Item `json:"item"`
	}
)

func New() *Postman {
	return &Postman{
		Item: make([]*Item, 0),
	}
}

func (o *Postman) Add(items ...*Item) {
	o.Item = append(o.Item, items...)
}

func (o *Postman) SetInfo(info *Info) {
	o.Info = info
}

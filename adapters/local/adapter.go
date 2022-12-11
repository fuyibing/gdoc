// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package local

type (
	Adapter interface {
	}

	adapter struct {
	}
)

func NewAdapter() Adapter {
	return (&adapter{}).init()
}

func (o *adapter) init() *adapter {
	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package postman

type (
	// Adapter
	//
	// Build postman.json to storage path, It is suitable for
	// Postman collect.
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

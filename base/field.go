// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package base

type (
	Field interface {
	}

	field struct {
	}
)

func NewField() Field {
	return (&field{}).init()
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *field) init() *field {
	return o
}

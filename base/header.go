// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

type (
	Header struct {
		Annotation  Annotation
		Key, Value  string
		Description string
	}
)

func NewHeader(a Annotation) *Header {
	return (&Header{Annotation: a}).init()
}

func (o *Header) init() *Header {
	o.Key = o.Annotation.GetValue(0)
	o.Value = o.Annotation.GetValue(1)
	o.Description = o.Annotation.GetValues(2)
	return o
}

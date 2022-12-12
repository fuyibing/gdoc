// author: wsfuyibing <websearch@163.com>
// date: 2022-12-12

package reflectors

type (
	FieldKind int
)

const (
	_ FieldKind = iota

	FieldJson
	FieldForm
	FieldUrl
)

// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package base

import (
	"strconv"
)

type (
	Error struct {
		Annotation Annotation

		Code                 int64
		Description, Message string
	}
)

func NewError(a Annotation) *Error {
	return (&Error{
		Annotation: a,
	}).init()
}

func (o *Error) init() *Error {
	if n, err := strconv.ParseInt(o.Annotation.GetFirst(), 10, 64); err == nil {
		o.Code = n
	}

	o.Message = o.Annotation.GetValue(1)
	o.Description = o.Annotation.GetValues(2)
	return o
}

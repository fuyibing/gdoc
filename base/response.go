// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

import (
	"fmt"
	"github.com/fuyibing/gdoc/conf"
	"regexp"
	"strings"
)

type (
	Response struct {
		Annotation Annotation
		Type       ResponseType

		Key       string
		Pkg, Name string
	}
)

func NewResponse(a Annotation, t ResponseType) *Response {
	return (&Response{Annotation: a, Type: t}).init()
}

func (o *Response) init() *Response {
	o.Key = strings.TrimPrefix(o.Annotation.GetFirst(), "/")

	if m := regexp.MustCompile(`^([^.]+)\.([^.]*)$`).FindStringSubmatch(o.Key); len(m) == 3 {
		o.Pkg = fmt.Sprintf("%s/%s", conf.Config.Module, m[1])
		o.Name = m[2]
	}

	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package base

import (
	"fmt"
	"github.com/fuyibing/gdoc/conf"
	"regexp"
	"strings"
)

type (
	Request struct {
		Annotation Annotation

		Key       string
		Pkg, Name string
	}
)

func NewRequest(a Annotation) *Request {
	return (&Request{Annotation: a}).init()
}

func (o *Request) init() *Request {
	o.Key = strings.TrimPrefix(o.Annotation.GetFirst(), "/")

	if m := regexp.MustCompile(`^([^.]+)\.([^.]*)$`).FindStringSubmatch(o.Key); len(m) == 3 {
		o.Pkg = fmt.Sprintf("%s/%s", conf.Config.Module, m[1])
		o.Name = m[2]
	}

	return o
}

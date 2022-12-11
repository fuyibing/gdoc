// author: wsfuyibing <websearch@163.com>
// date: 2022-12-10

package base

import (
	"github.com/fuyibing/gdoc/config"
)

type (
	// Request
	// defined by annotation.
	Request interface {
		GetKey() string
		GetLine() int
		GetName() string
		GetPackage() string
	}

	request struct {
		Line    int
		Key     string // eg. app/logics/user.LoginRequest
		Name    string // eg. LoginRequest
		Package string // eg. app/logics/user
	}
)

func NewRequest(n int, k string) Request {
	return (&request{Line: n, Key: k}).init()
}

func (o *request) GetKey() string     { return o.Key }
func (o *request) GetLine() int       { return o.Line }
func (o *request) GetName() string    { return o.Name }
func (o *request) GetPackage() string { return o.Package }

func (o *request) init() *request {
	if m := config.Regex.GetStructWithPackage().FindStringSubmatch(o.Key); len(m) == 3 {
		o.Name = m[2]
		o.Package = m[1]
	}
	return o
}

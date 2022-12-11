// author: wsfuyibing <websearch@163.com>
// date: 2022-12-10

package base

import (
	"github.com/fuyibing/gdoc/config"
)

type (
	// Response
	// defined by annotation.
	Response interface {
		GetKey() string
		GetLine() int
		GetName() string
		GetPackage() string
		GetType() ResponseType
	}

	ResponseType int

	response struct {
		Line int
		Type ResponseType

		Key     string // eg. app/logics/user.LoginResponse
		Name    string // eg. LoginResponse
		Package string // eg. app/logics/user
	}
)

const (
	_ ResponseType = iota

	ResponseData
	ResponseList
	ResponsePaging
)

func NewResponse(t ResponseType, n int, k string) Response {
	return (&response{
		Line: n,
		Key:  k,
		Type: t,
	}).init()
}

func (o *response) GetKey() string        { return o.Key }
func (o *response) GetLine() int          { return o.Line }
func (o *response) GetName() string       { return o.Name }
func (o *response) GetPackage() string    { return o.Package }
func (o *response) GetType() ResponseType { return o.Type }

func (o *response) init() *response {
	if m := config.Regex.GetStructWithPackage().FindStringSubmatch(o.Key); len(m) == 3 {
		o.Name = m[2]
		o.Package = m[1]
	}
	return o
}

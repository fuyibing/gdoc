// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package base

import (
	"fmt"
	"github.com/fuyibing/gdoc/config"
	"strings"
)

type (
	Controller interface {
		GetComment() Comment
		GetMethod(key string) Method
		GetMethods() map[string]Method
		GetName() string
		GetPackage() string
		GetPrefix() string
		SetMethod(key string, method Method)
		SetPrefix(prefix string)
	}

	controller struct {
		Comment Comment

		Name    string // eg. UserController
		Package string // eg. app/controllers/user
		Prefix  string // eg. /user

		methodList   []string
		methodMapper map[string]Method
	}
)

func NewController(name, prefix string) Controller {
	return (&controller{
		Comment: NewComment(name),

		Name:    name,
		Prefix:  prefix,
		Package: strings.TrimPrefix(fmt.Sprintf("%s%s", config.Path.GetController(), prefix), "/"),
	}).init()
}

// /////////////////////////////////////////////////////////////
// Interface method
// /////////////////////////////////////////////////////////////

func (o *controller) GetComment() Comment                 { return o.Comment }
func (o *controller) GetMethod(key string) Method         { return o.getMethod(key) }
func (o *controller) GetMethods() map[string]Method       { return o.methodMapper }
func (o *controller) GetName() string                     { return o.Name }
func (o *controller) GetPackage() string                  { return o.Package }
func (o *controller) GetPrefix() string                   { return o.Prefix }
func (o *controller) SetMethod(key string, method Method) { o.methodMapper[key] = method }
func (o *controller) SetPrefix(prefix string)             { o.Prefix = prefix }

// /////////////////////////////////////////////////////////////
// Access method
// /////////////////////////////////////////////////////////////

func (o *controller) getMethod(key string) Method {
	if m, ok := o.methodMapper[key]; ok {
		return m
	}
	return nil
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *controller) init() *controller {
	o.methodList = make([]string, 0)
	o.methodMapper = make(map[string]Method)
	return o
}

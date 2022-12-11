// author: wsfuyibing <websearch@163.com>
// date: 2022-12-06

package reflectors

import (
	"fmt"
	"regexp"
)

var (
	ErrAnonymousPointerNotAllowed = fmt.Errorf("pointer anonymous not allow")
	ErrAnonymousUnknown           = fmt.Errorf("unknown anonymous not allow")
	ErrUnknownType                = fmt.Errorf("unknown field type")
)

var (
	RegexpFieldValidate = regexp.MustCompile(`required`)
)

const (
	DefaultBoolValue      = false
	DefaultFloatValue     = 0
	DefaultInterfaceValue = "*"
	DefaultIntValue       = 0
	DefaultMapValue       = "{}"
	DefaultStringValue    = ""
)

const (
	TagDescription = "description"
	TagExec        = "exec"
	TagIgnored     = "-"
	TagJson        = "json"
	TagLabel       = "label"
	TagMock        = "mock"
	TagValidate    = "validate"
)

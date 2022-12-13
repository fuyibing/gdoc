// author: wsfuyibing <websearch@163.com>
// date: 2022-12-13

package local

import (
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"strings"
	"time"
)

type (
	MethodTpl struct {
		adapter  *adapter
		method   base.Method
		contents []string
	}
)

func NewMethodTpl(a *adapter, m base.Method) *MethodTpl {
	return (&MethodTpl{
		adapter: a,
		method:  m,
	}).init()
}

// Save
// 保存文档.
func (o *MethodTpl) Save() error {
	for _, g := range [][]func(){
		// 1. Header.
		{
			o.title,
			o.description,
			o.code,
		},

		// 2. ErrorCode
		// 2. Request
		// 3. Response
	} {
		for _, call := range g {
			call()
		}
	}

	// 写入目标.
	return o.adapter.manager.SaveFile(
		o.adapter.genMethodFilename(o.method.GetRequestMethod(), o.method.GetRequestUrl(), o.method.GetController().GetPrefix()),
		strings.Join(o.contents, "\n\n"),
	)
}

func (o *MethodTpl) init() *MethodTpl {
	o.contents = make([]string, 0)
	return o
}

// /////////////////////////////////////////////////////////////
// Part 1
// /////////////////////////////////////////////////////////////

func (o *MethodTpl) code() {
	var (
		comma = "//\n"
		list  = []string{"```go"}
	)

	// 位置.
	list = append(list,
		fmt.Sprintf("// Date: %v", time.Now().Format("2006-01-02 15:04")),
		fmt.Sprintf("// Line: %d", o.method.GetComment().GetLine()),
		fmt.Sprintf("// File: %v", o.method.GetComment().GetPath()),
		"",
	)

	// 注释.
	list = append(list,
		fmt.Sprintf("// %v", o.method.GetName()),
		fmt.Sprintf("// %v.", o.method.GetComment().GetTitle()),
	)

	// 入参.
	if r := o.method.GetComment().GetRequest(); r != nil {
		list = append(list, fmt.Sprintf("%s// @Request(%s)", comma, r.GetKey()))
		comma = ""
	}

	// 出参.
	if rs := o.method.GetComment().GetResponses(); len(rs) > 0 {
		for _, r := range rs {
			list = append(list, fmt.Sprintf("%s// @Response(%s)", comma, r.GetKey()))
			comma = ""
		}
	}

	// 源码.
	list = append(list, o.method.GetComment().GetCode(), "    // ...", "}")
	list = append(list, "```")

	// 加入.
	o.contents = append(o.contents, strings.Join(list, "\n"))
}

func (o *MethodTpl) description() {
	if desc := o.method.GetComment().GetDescription(); desc != "" {
		o.contents = append(o.contents, fmt.Sprintf(
			"> %s",
			strings.ReplaceAll(desc, "\n", "\n> "),
		))
	}
}

func (o *MethodTpl) title() {
	o.contents = append(o.contents, fmt.Sprintf("# %s", o.method.GetComment().GetTitle()))
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-13

package local

import (
	"encoding/json"
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
		// 前置.
		{
			o.title,
			o.info,
			o.description,
			o.code,
		},

		// 入参.
		{
			o.errors,
			o.headers,
			o.request,
		},

		// 出参.
		{
			o.response,
		},
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
	if o.method.GetName() == o.method.GetComment().GetTitle() {
		list = append(list,
			fmt.Sprintf("// %v .", o.method.GetName()),
		)
	} else {
		list = append(list,
			fmt.Sprintf("// %v", o.method.GetName()),
			fmt.Sprintf("// %v.", o.method.GetComment().GetTitle()),
		)
	}

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
	o.contents = append(o.contents,
		fmt.Sprintf("# %s", o.method.GetComment().GetTitle()),
	)
}

// /////////////////////////////////////////////////////////////
// Part 2
// /////////////////////////////////////////////////////////////

func (o *MethodTpl) info() {
	var (
		route = fmt.Sprintf("%s%s", o.method.GetController().GetPrefix(), o.method.GetRequestUrl())
	)

	if route == "" {
		route = "/"
	}

	o.contents = append(o.contents,
		"----",
		fmt.Sprintf("`%s` `%s` `%s`",
			o.method.GetRequestMethod(),
			route,
			o.method.GetContentType(),
		),
		"----",
	)
}

// /////////////////////////////////////////////////////////////
// Part 3
// /////////////////////////////////////////////////////////////

func (o *MethodTpl) errors() {
	var (
		list = make([]string, 0)
		errs = o.method.GetComment().GetErrors()
	)

	// 无值.
	if len(errs) == 0 {
		o.contents = append(o.contents, "----",
			fmt.Sprintf("### 错误码"),
			fmt.Sprintf("> 空"),
		)
		return
	}

	// 表头.
	list = append(list,
		fmt.Sprintf("| 编码 | 错误 | 描述 |"),
		fmt.Sprintf("| :---- | :---- | :---- |"),
	)

	// 表体.
	for _, err := range errs {
		list = append(list,
			fmt.Sprintf("| %v | %v | %v |",
				err.Code,
				err.Message,
				err.Description,
			),
		)
	}

	// 加入.
	o.contents = append(o.contents,
		"----",
		fmt.Sprintf("### 错误码"),
		strings.Join(list, "\n"),
	)
}

func (o *MethodTpl) headers() {
	var (
		list    = make([]string, 0)
		headers = o.method.GetComment().GetHeaders()
	)

	// 无值.
	if len(headers) == 0 {
		o.contents = append(o.contents, "----",
			fmt.Sprintf("### 请求头"),
			fmt.Sprintf("> 空"),
		)
		return
	}

	// 表头.
	list = append(list,
		fmt.Sprintf("| 键名 | 键值 | 描述 |"),
		fmt.Sprintf("| :---- | :---- | :---- |"),
	)

	// 表体.
	for _, header := range headers {
		list = append(list,
			fmt.Sprintf("| %v | %v | %v |",
				header.Key,
				header.Value,
				header.Description,
			),
		)
	}

	// 加入.
	o.contents = append(o.contents,
		"----",
		fmt.Sprintf("### 请求头"),
		strings.Join(list, "\n"),
	)
}

func (o *MethodTpl) request() {
	if r := o.method.GetComment().GetRequest(); r != nil {
		o.contents = append(o.contents,
			"----",
			"### 入参",
			fmt.Sprintf("> %s", r.GetKey()),
		)

		o.render(r.GetKey(), true)
	}
}

func (o *MethodTpl) response() {
	ls := o.method.GetComment().GetResponses()
	if n := len(ls); n > 0 {
		for i, r := range ls {
			o.contents = append(o.contents,
				"----",
				fmt.Sprintf("### 出叁: %d/%d", i+1, n),
				fmt.Sprintf("> %s", r.GetKey()),
			)

			o.render(r.GetKey(), false)
		}
	}
}

func (o *MethodTpl) render(key string, input bool) {
	o.renderTable(key, input)
	o.renderCode(key)
}

// 渲染编码.
func (o *MethodTpl) renderCode(key string) {
	if v := o.adapter.manager.LoadJsonCode(key); v != nil {
		if buf, err := json.MarshalIndent(v, "", "    "); err == nil {
			o.contents = append(o.contents, "> Code", strings.Join([]string{
				"```json", string(buf), "```",
			}, "\n"))
		}
	}
}

// 渲染表格.
func (o *MethodTpl) renderTable(key string, input bool) {
	var (
		items = o.adapter.manager.LoadJsonItems(key)
		table = make([]string, 0)
	)

	// 字段.
	if len(items) == 0 {
		return
	}

	// 表头.
	if input {
		table = append(table,
			fmt.Sprintf("| 字段 | 类型 | 必需 | 条件 | 说明 | 参考值 |"),
			fmt.Sprintf("| :---- | :---- | :----: | :---- | :---- | :---- |"),
		)
	} else {
		table = append(table,
			fmt.Sprintf("| 字段 | 类型 | 说明 | 参考值 |"),
			fmt.Sprintf("| :---- | :---- | :---- | :---- |"),
		)
	}

	// 表体.
	if vs := o.renderTbody(0, items, input); len(vs) > 0 {
		table = append(table, vs...)
	}

	// 加入.
	o.contents = append(o.contents, strings.Join(table, "\n"))
}

func (o *MethodTpl) renderTbody(offset int, items []*base.Item, input bool) (results []string) {
	results = make([]string, 0)

	for _, item := range items {
		// 字段行.
		if input {
			// 入参.
			results = append(results, fmt.Sprintf(
				"| %v | `%v` | %v | %v | %v | %v |",
				o.generateKey(offset, item),
				o.generateType(item),
				o.generateRequired(item),
				item.Condition,
				o.generateLabel(item),
				o.generateValue(item),
			))
		} else {
			// 出参.
			results = append(results, fmt.Sprintf(
				"| %v | `%v` | %v | %v |",
				o.generateKey(offset, item),
				o.generateType(item),
				o.generateLabel(item),
				o.generateValue(item),
			))
		}

		// 子嵌套.
		if item.Children != nil && len(item.Children) > 0 {
			results = append(results, o.renderTbody(offset+1, item.Children, input)...)
		}
	}

	return
}

func (o *MethodTpl) generateKey(offset int, item *base.Item) string {
	return fmt.Sprintf("%s %s",
		strings.Repeat("　 ", offset),
		item.Key,
	)
}

func (o *MethodTpl) generateLabel(item *base.Item) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.Join([]string{
		item.Label, item.Description,
	}, "\n")), "\n", "<br />")
}

func (o *MethodTpl) generateType(item *base.Item) string {
	if item.Array {
		return fmt.Sprintf("[] %s", item.Type)
	}
	return item.Type
}

func (o *MethodTpl) generateRequired(item *base.Item) string {
	if item.Required {
		return "Y"
	}
	return " "
}

func (o *MethodTpl) generateValue(item *base.Item) string {
	return fmt.Sprintf("%v", item.Value)
}

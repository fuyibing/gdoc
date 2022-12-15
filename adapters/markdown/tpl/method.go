// author: wsfuyibing <websearch@163.com>
// date: 2022-12-15

package tpl

import (
	"fmt"
	"github.com/fuyibing/gdoc/adapters/markdown/i18n"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/conf"
	"strings"
)

type (
	Method struct {
		mapping   base.Mapping
		method    base.Method
		templates []string
	}
)

func NewMethod(mapping base.Mapping, method base.Method) *Method {
	return (&Method{
		mapping:   mapping,
		method:    method,
		templates: make([]string, 0),
	}).init()
}

func (o *Method) Save() {
	for _, call := range []func(){
		o.title,
		o.info,
		o.description,

		o.code,

		o.request,
		o.responses,

		o.save,
	} {
		call()
	}
}

func (o *Method) code() {
	var (
		list = make([]string, 0)
	)

	list = append(list, "```go",
		fmt.Sprintf("// %s: <%v:%d>", i18n.Lang("Source Path"), o.method.GetComment().Path, o.method.GetComment().Line),
		fmt.Sprintf("// %s: <%v.%v>", i18n.Lang("Source Package"), o.method.GetComment().Pkg, o.method.GetController().GetComment().Name),
		"",
	)

	// Head comment start with name.
	if o.method.GetComment().Name == o.method.GetComment().GetTitle() {
		list = append(list,
			fmt.Sprintf("// %s .", o.method.GetComment().Name),
		)
	} else {
		list = append(list,
			fmt.Sprintf("// %s", o.method.GetComment().Name),
			fmt.Sprintf("// %s.", o.method.GetComment().GetTitle()),
		)
	}

	// Header.
	if len(o.method.GetComment().Headers) > 0 {
		list = append(list, "//")
		for _, r := range o.method.GetComment().Headers {
			list = append(list, fmt.Sprintf(
				"// @%s(\"%s\", \"%s\", \"%v\")",
				r.Annotation.GetName(),
				r.Key,
				r.Value,
				r.Description,
			))
		}
	}

	// Request.
	if r := o.method.GetComment().Request; r != nil {
		list = append(list, "//", fmt.Sprintf(
			"// @%s(%s)",
			r.Annotation.GetName(),
			r.Key,
		))
	}

	// Response.
	if len(o.method.GetComment().Responses) > 0 {
		list = append(list, "//")
		for _, r := range o.method.GetComment().Responses {
			list = append(list, fmt.Sprintf(
				"// @%s(%s)",
				r.Annotation.GetName(),
				r.Key,
			))
		}
	}

	// Error codes.
	if len(o.method.GetComment().Errors) > 0 {
		list = append(list, "//")
		for _, r := range o.method.GetComment().Errors {
			list = append(list, fmt.Sprintf(
				"// @%s(%v, \"%s\", \"%v\")",
				r.Annotation.GetName(),
				r.Code,
				r.Message,
				r.Description,
			))
		}
	}

	// Source code.
	list = append(list,
		fmt.Sprintf("%s",
			strings.TrimSuffix(o.method.GetComment().Code, "}"),
		), "}",
	)

	list = append(list, "```")
	o.templates = append(o.templates, strings.Join(list, "\n"))
}

func (o *Method) description() {
	if str := o.method.GetComment().GetDescription(); str != "" {
		list := make([]string, 0)
		for _, s := range strings.Split(str, "\n") {
			if s = strings.TrimSpace(s); s != "" {
				list = append(list, fmt.Sprintf("> %s", s))
			}
		}
		o.templates = append(o.templates, strings.Join(list, "\n"))
	}
}

func (o *Method) title() {
	o.templates = append(o.templates,
		fmt.Sprintf("# %s", o.method.GetComment().GetTitle()),
	)
}

func (o *Method) info() {
	list := []string{
		fmt.Sprintf("|  |  |"),
		fmt.Sprintf("| ----: | :---- |"),
		fmt.Sprintf("| %s | `%s` |", i18n.Lang("Request Domain"), conf.Config.Deploy.Full()),
		fmt.Sprintf("| %s | `%s` `%s` |", i18n.Lang("Request URL"), o.method.GetRequestMethod(), o.method.GetRequestUrl()),
		fmt.Sprintf("| %s | `%s`| ", i18n.Lang("Content Type"), o.method.GetContentType()),
		fmt.Sprintf("| %s | `%s` |", i18n.Lang("Last Updated"), o.mapping.GetLastUpdated()),
	}

	o.templates = append(o.templates, strings.Join(list, "\n"))
}

func (o *Method) init() *Method {
	return o
}

// Build request params.
func (o *Method) request() {
	var (
		list = make([]string, 0)
		req  = o.method.GetComment().Request
	)

	// Not registered by annotation.
	if req == nil {
		return
	}

	// Request fields.
	if ls := o.renderTable(req.Key, true); len(ls) > 0 {
		list = append(list, ls...)
	}

	// items := o.mapping.LoadTmpItem(req.Key)
	// if len(items) > 0 {
	// 	o.requestFields(items)
	// }

	o.templates = append(o.templates,
		"----",
		fmt.Sprintf("### %s", i18n.Lang("Request Params")),
		fmt.Sprintf("> `%s`.`%s`", req.Pkg, req.Name),
		strings.Join(list, "\n"),
	)
}

// Build response lists.
func (o *Method) response(i int, r *base.Response) {
	var (
		list = make([]string, 0)
	)

	// Request fields.
	if ls := o.renderTable(r.Key, false); len(ls) > 0 {
		list = append(list, ls...)
	}

	o.templates = append(o.templates,
		"----",
		fmt.Sprintf("### %s %d", i18n.Lang("Response Params"), i+1),
		fmt.Sprintf("> `%s`.`%s`", r.Pkg, r.Name),
		strings.Join(list, "\n"),
	)
}

func (o *Method) responses() {
	for i, r := range o.method.GetComment().Responses {
		o.response(i, r)
	}
}

func (o *Method) save() {
	conf.Path.SavePath(
		fmt.Sprintf("%s%s%s",
			conf.Path.GetBasePath(),
			conf.Path.GetDocumentPath(),
			linkName(o.method),
		),
		strings.Join(o.templates, "\n\n"),
	)
}

// /////////////////////////////////////////////////////////////
// Render fields.
// /////////////////////////////////////////////////////////////

func (o *Method) renderCellDesc(item *base.Item) string {
	return strings.ReplaceAll(
		strings.TrimSpace(
			strings.Join([]string{
				item.Label,
				item.Description,
			}, "\n"),
		), "\n", "<br />",
	)
}

func (o *Method) renderCellKey(n int, item *base.Item) string {
	return fmt.Sprintf("%s %s",
		strings.Repeat("　　", n),
		item.Key,
	)
}

func (o *Method) renderCellRequired(item *base.Item) string {
	if item.Required {
		return "Y"
	}
	return ""
}

func (o *Method) renderCellType(item *base.Item) string {
	if item.Array {
		return fmt.Sprintf("[]%s", item.Type)
	}
	return item.Type
}

func (o *Method) renderTable(key string, input bool) []string {
	if items := o.mapping.LoadTmpItem(key); len(items) > 0 {
		if ls := o.renderTableBody(input, 0, items); len(ls) > 0 {
			list := o.renderTableHeader(input)
			list = append(list, ls...)
			return list
		}
	}
	return nil
}

func (o *Method) renderTableBody(input bool, n int, items []*base.Item) []string {
	list := make([]string, 0)

	for _, item := range items {
		if input {
			list = append(list, fmt.Sprintf(
				"| %v | `%v` | %v | %v | %v |",
				o.renderCellKey(n, item),
				o.renderCellType(item),
				o.renderCellRequired(item),
				item.Condition,
				o.renderCellDesc(item),
			))
		} else {
			list = append(list, fmt.Sprintf(
				"| %v | `%v` | %v |",
				o.renderCellKey(n, item),
				o.renderCellType(item),
				o.renderCellDesc(item),
			))
		}

		if len(item.Children) > 0 {
			if children := o.renderTableBody(input, n+1, item.Children); len(children) > 0 {
				list = append(list, children...)
			}
		}
	}

	return list
}

func (o *Method) renderTableHeader(input bool) []string {
	if input {
		return []string{
			fmt.Sprintf("| %v | %v | %v | %v | %v |",
				i18n.Lang("Field Name"),
				i18n.Lang("Field Type"),
				i18n.Lang("Required"),
				i18n.Lang("Validation"),
				i18n.Lang("Description"),
			), fmt.Sprintf(
				"| ---- | ---- | :----: | :---- | ---- |",
			),
		}
	}

	return []string{
		fmt.Sprintf("| %v | %v | %v |",
			i18n.Lang("Field Name"),
			i18n.Lang("Field Type"),
			i18n.Lang("Description"),
		), fmt.Sprintf(
			"| ---- | ---- | ---- |",
		),
	}
}

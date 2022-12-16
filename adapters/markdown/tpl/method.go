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
		mapping base.Mapping
		method  base.Method

		nl, sep   string
		pre       string
		templates []string
	}
)

// NewMethod
// create and return instance.
func NewMethod(mapping base.Mapping, method base.Method) *Method {
	return (&Method{
		mapping: mapping,
		method:  method,
	}).init()
}

// Save
// markdown file to specified directory.
func (o *Method) Save() {
	for _, call := range []func(){
		o.HeadTitle,
		o.ApiDesc,
		o.ApiInfo,
		o.ApiCode,

		o.RequestHeader,
		o.RequestParam,
		o.ErrorCode,
		o.ResponseParams,

		o.save,
	} {
		call()
	}
}

// /////////////////////////////////////////////////////////////
// Append callables methods
// /////////////////////////////////////////////////////////////

// ApiCode
// add example code.
func (o *Method) ApiCode() {
	var (
		pre  = "  "
		list = []string{
			fmt.Sprintf("%s```go", pre),
			fmt.Sprintf("%s// %s", pre, o.method.GetController().GetComment().Name),
			fmt.Sprintf("%s// %s.", pre, o.method.GetController().GetComment().GetTitle()),
			fmt.Sprintf("%stype %s struct {", pre, strings.TrimSuffix(o.method.GetController().GetComment().Name, "}")),
			fmt.Sprintf("%s}", pre),
			fmt.Sprintf("%s", pre),
			fmt.Sprintf("%s// %s", pre, o.method.GetComment().Name),
			fmt.Sprintf("%s// %s.", pre, o.method.GetComment().GetTitle()),
			fmt.Sprintf("%s%s", pre, strings.TrimSuffix(o.method.GetComment().Code, "}")),
			fmt.Sprintf("%s}", pre),
			fmt.Sprintf("%s```", pre),
		}
	)

	// Append to template.
	o.templates = append(o.templates, strings.Join(list, o.nl))
}

// ApiDesc
// add api description.
func (o *Method) ApiDesc() {
	if str := o.method.GetComment().GetDescription(); str != "" {
		o.templates = append(o.templates, "----", str)
	}
}

// ApiInfo
// add request route and content type info.
func (o *Method) ApiInfo() {
	o.templates = append(o.templates, "----", strings.Join([]string{
		fmt.Sprintf(
			"* %v: `%s`",
			i18n.Lang("Content Type"),
			o.method.GetContentType(),
		), fmt.Sprintf(
			"* %v: `%v` `%v`",
			i18n.Lang("Route"),
			o.method.GetRequestMethod(),
			o.method.GetRequestUrl(),
		), fmt.Sprintf(
			"* %v: `%v`",
			i18n.Lang("Deploy"),
			conf.Config.Deploy.Full(),
		), fmt.Sprintf(
			"* %v: `%v`",
			i18n.Lang("Updated"),
			o.mapping.GetLastUpdated(),
		), fmt.Sprintf("* %s: `<%s.%s> %s`",
			i18n.Lang("Struct"),
			o.method.GetController().GetComment().Pkg,
			o.method.GetController().GetComment().Name,
			o.method.GetComment().Name,
		),
	}, o.nl))

}

// ErrorCode
// add error codes.
func (o *Method) ErrorCode() {
	if len(o.method.GetComment().Errors) < 1 {
		o.templates = append(o.templates,
			fmt.Sprintf("### %v", i18n.Lang("Error Codes")),
		)
		return
	}

	list := []string{
		fmt.Sprintf(
			"| %v | %v | %v |",
			i18n.Lang("Code"),
			i18n.Lang("Value"),
			i18n.Lang("Description"),
		),
		fmt.Sprintf("| ---- | ---- | ---- |"),
	}

	for _, x := range o.method.GetComment().Errors {
		list = append(list, fmt.Sprintf(
			"| %v | %v | %v |",
			x.Code, x.Message, x.Description,
		))
	}

	o.templates = append(o.templates,
		fmt.Sprintf("### %v", i18n.Lang("Error Codes")),
		strings.Join(list, o.nl),
	)
}

// HeadTitle
// add H1 tag.
func (o *Method) HeadTitle() {
	o.templates = append(o.templates, fmt.Sprintf("# %s",
		o.method.GetComment().GetTitle(),
	))
}

// RequestHeader
// add request header params.
func (o *Method) RequestHeader() {
	// H3
	// for request params.
	o.templates = append(o.templates, fmt.Sprintf("### %v", i18n.Lang("Request Headers")))

	// Undefined.
	if len(o.method.GetComment().Headers) < 1 {
		return
	}

	list := []string{
		fmt.Sprintf(
			"| %v | %v | %v |",
			i18n.Lang("Key"),
			i18n.Lang("Value"),
			i18n.Lang("Description"),
		),
		fmt.Sprintf("| ---- | ---- | ---- |"),
	}

	for _, x := range o.method.GetComment().Headers {
		list = append(list, fmt.Sprintf(
			"| %v | %v | %v |",
			x.Key, x.Value, x.Description,
		))
	}

	o.templates = append(o.templates,
		fmt.Sprintf("### %v", i18n.Lang("Request Headers")),
		strings.Join(list, o.nl),
	)
}

// RequestParam
// add request param.
func (o *Method) RequestParam() {
	// H3
	// title for request annotation.
	//
	//   # Request Params
	o.templates = append(o.templates, fmt.Sprintf("### %v", i18n.Lang("Request Params")))

	// Return
	// if request undefined.
	if o.method.GetComment().Request == nil {
		return
	}

	// Request
	// struct info.
	//
	//   * Struct: `LoginRequest`
	//   * Package: `mod/app/logics/user`
	o.templates = append(o.templates,
		strings.Join([]string{
			fmt.Sprintf("* %s: `%s`",
				i18n.Lang("Struct"),
				o.method.GetComment().Request.Name,
			), fmt.Sprintf("* %s: `%s`",
				i18n.Lang("Package"),
				o.method.GetComment().Request.Pkg,
			),
		}, o.nl),
	)

	// Request
	// definition fields.
	if list := o.table(true, o.method.GetComment().Request.Key); len(list) > 0 {
		o.templates = append(o.templates, strings.Join(list, o.nl))
	}

	// Request
	// json code.
	//
	//   *Example Code*
	//
	//   ```json
	//   {
	//       "key": "value"
	//   }
	//   ```
	if code := string(o.mapping.LoadTmpCode(o.method.GetComment().Request.Key)); code != "" {
		// json code.
		code = strings.ReplaceAll(code,
			fmt.Sprintf("\n"),
			fmt.Sprintf("\n%s", o.pre),
		)

		// append to templates.
		o.templates = append(o.templates,
			// Catalog.
			fmt.Sprintf("%s*%v*: ", o.pre, i18n.Lang("Example Code")),

			// JSON Code.
			strings.Join([]string{
				fmt.Sprintf("%s```json", o.pre),
				fmt.Sprintf("%s%s", o.pre, code),
				fmt.Sprintf("%s```", o.pre),
			}, o.nl),
		)
	}
}

// ResponseParam
// add response param.
func (o *Method) ResponseParam(i int, r *base.Response) {
	o.templates = append(o.templates, fmt.Sprintf("### %v # %d", i18n.Lang("Response Params"), i+1))

	// Request
	// struct info.
	//
	//   * Struct: `LoginRequest`
	//   * Package: `mod/app/logics/user`
	o.templates = append(o.templates,
		strings.Join([]string{
			fmt.Sprintf("* %s: `%s`",
				i18n.Lang("Struct"),
				r.Name,
			), fmt.Sprintf("* %s: `%s`",
				i18n.Lang("Package"),
				r.Pkg,
			),
		}, o.nl),
	)

	// Response
	// definition fields.
	if list := o.table(false, r.Key); len(list) > 0 {
		o.templates = append(o.templates, strings.Join(list, o.nl))
	}

	// Response
	// json code.
	//
	//   *Example Code*
	//
	//   ```json
	//   {
	//       "key": "value"
	//   }
	//   ```
	if code := r.Type.Render(o.mapping.LoadTmpCode(r.Key)); code != "" {
		// json code.
		code = strings.ReplaceAll(code,
			fmt.Sprintf("\n"),
			fmt.Sprintf("\n%s", o.pre),
		)

		// append to templates.
		o.templates = append(o.templates,
			// Catalog.
			fmt.Sprintf("%s*%v*: ", o.pre, i18n.Lang("Example Code")),

			// JSON Code.
			strings.Join([]string{
				fmt.Sprintf("%s```json", o.pre),
				fmt.Sprintf("%s%s", o.pre, code),
				fmt.Sprintf("%s```", o.pre),
			}, o.nl),
		)
	}
}

// ResponseParams
// iterate response definitions.
func (o *Method) ResponseParams() {
	if len(o.method.GetComment().Responses) < 1 {
		o.templates = append(o.templates, fmt.Sprintf("### %v", i18n.Lang("Response Params")))
		return
	}

	for i, r := range o.method.GetComment().Responses {
		o.ResponseParam(i, r)
	}
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Method) renderCondition(x *base.Item) string {
	return x.Condition
}

func (o *Method) renderDescription(x *base.Item) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.TrimSpace(strings.Join(
				[]string{x.Label, x.Description}, "\n",
			)), "\n", ", ",
		), "|", "｜",
	)
}

func (o *Method) renderExample(x *base.Item) string {
	return fmt.Sprintf("%v", x.Value)
}

func (o *Method) renderKey(n int, x *base.Item) string {
	return fmt.Sprintf("%s%s",
		strings.Repeat("　 ", n),
		x.Key,
	)
}

func (o *Method) renderRequired(x *base.Item) string {
	if x.Required {
		return "`Y`"
	}
	return " "
}

func (o *Method) renderType(x *base.Item) string {
	if x.Array {
		return fmt.Sprintf("`[]` `%s`", x.Type)
	}
	return fmt.Sprintf("`%s`", x.Type)
}

// Table collect.
func (o *Method) table(input bool, key string) []string {
	var (
		list  = make([]string, 0)
		items = o.mapping.LoadTmpItem(key)
	)

	// Return empty
	// if source file not found or fields not exported.
	if len(items) < 1 {
		return list
	}

	// Append
	// head and body.
	list = append(list, o.thead(input)...)
	list = append(list, o.tbody(input, 0, items)...)
	return list
}

// Table body collect.
func (o *Method) tbody(input bool, n int, items []*base.Item) []string {
	list := make([]string, 0)

	// Iterate
	// items.
	for _, item := range items {
		if input {
			list = append(list, fmt.Sprintf("%s| %v | %v | %v | %v | %v | %v |",
				o.pre,
				o.renderKey(n, item),
				o.renderType(item),
				o.renderRequired(item),
				o.renderCondition(item),
				o.renderDescription(item),
				o.renderExample(item),
			))
		} else {
			list = append(list, fmt.Sprintf("%s| %v | %v | %v |",
				o.pre,
				o.renderKey(n, item),
				o.renderType(item),
				o.renderDescription(item),
			))
		}

		// Recursion
		// on child items.
		if len(item.Children) > 0 {
			list = append(list, o.tbody(input, n+1, item.Children)...)
		}
	}
	return list
}

// Table head collect.
func (o *Method) thead(input bool) []string {
	// Input
	// for request.
	//
	//   | Field | Type | Required | Condition | Description | Example |
	//   | ----- | ---- | -------- | --------- | ----------- | ------- |
	if input {
		return []string{
			fmt.Sprintf("%s| %v | %v | %v | %v | %v | %v |",
				o.pre,
				i18n.Lang("Field"),
				i18n.Lang("Type"),
				i18n.Lang("Required"),
				i18n.Lang("Condition"),
				i18n.Lang("Description"),
				i18n.Lang("Example"),
			), fmt.Sprintf("%s| ---- | ---- | :----: | ---- | ---- | ---- |",
				o.pre,
			),
		}
	}

	// Simple
	// for response.
	//
	//   | Field | Type | Description |
	//   | ----- | ---- | ----------- |
	return []string{
		fmt.Sprintf("%s| %v | %v | %v |",
			o.pre,
			i18n.Lang("Field"),
			i18n.Lang("Type"),
			i18n.Lang("Description"),
		), fmt.Sprintf("%s| ---- | ---- | ---- |",
			o.pre,
		),
	}
}

// /////////////////////////////////////////////////////////////
// Initialization
// /////////////////////////////////////////////////////////////

func (o *Method) init() *Method {
	o.nl = "\n"
	o.sep = "\n\n"
	o.pre = "  "
	o.templates = make([]string, 0)
	return o
}

func (o *Method) save() {
	// Append source info
	// to markdown files end.
	o.templates = append(o.templates, "----",
		strings.Join([]string{
			fmt.Sprintf("* %s: `%v`", i18n.Lang("Line"), o.method.GetComment().Line),
			fmt.Sprintf("* %s: `%v`", i18n.Lang("Path"), o.method.GetComment().Path),
			fmt.Sprintf("* %s: `%s`.`%s`", i18n.Lang("Package"), o.method.GetComment().Pkg, o.method.GetComment().Name),
		}, o.nl),
	)

	// Call save handler
	// to specified file.
	conf.Path.SavePath(
		fmt.Sprintf("%s%s%s",
			conf.Path.GetBasePath(),
			conf.Path.GetDocumentPath(),
			linkName(o.method),
		), strings.Join(
			o.templates,
			o.sep,
		),
	)
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package base

import (
	"github.com/fuyibing/gdoc/config"
	"strings"
)

type (
	Comment interface {
		AddText(line int, text string)

		GetCode() string
		GetLine() int
		GetPath() string

		GetTitle() (s string)
		GetDescription() (s string)
		GetRequest() Request
		GetResponses() []Response
		IsIgnored() bool
		SetSource(line int, code, path string)
		SetRequest(n int, s string)
		SetResponse(t ResponseType, n int, s string)
	}

	comment struct {
		code, path string
		line       int

		name       string
		request    Request
		responses  []Response
		textLength int32
		textList   []string
	}
)

func NewComment(name string) Comment {
	return (&comment{name: name}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *comment) AddText(line int, text string) {
	// Remove
	// end string with dot.
	text = config.Regex.GetCommentEnd().ReplaceAllString(text, "")

	// First
	// comment start with name.
	if o.textLength == 0 {
		if s := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(text), o.name)); s != "" {
			text = s
		} else {
			return
		}
	}

	// Add to list.
	if text != "" {
		o.textLength++
		o.textList = append(o.textList, text)
	}
}

func (o *comment) GetCode() string { return o.code }
func (o *comment) GetLine() int    { return o.line }
func (o *comment) GetPath() string { return o.path }

func (o *comment) GetTitle() string {
	if o.textLength > 0 {
		return strings.TrimSpace(o.textList[0])
	}
	return o.name
}

func (o *comment) GetDescription() (s string) {
	if o.textLength > 1 {
		s = strings.TrimSpace(strings.Join(o.textList[1:], "\n"))
	}
	return
}

func (o *comment) GetRequest() Request { return o.request }

func (o *comment) GetResponses() []Response { return o.responses }

func (o *comment) IsIgnored() bool { return false }

func (o *comment) SetRequest(n int, s string) { o.request = NewRequest(n, s) }

func (o *comment) SetResponse(t ResponseType, n int, s string) {
	o.responses = append(o.responses, NewResponse(t, n, s))
}

func (o *comment) SetSource(line int, code, path string) {
	o.line = line
	o.code = code
	o.path = path
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *comment) init() *comment {
	o.textList = make([]string, 0)
	o.responses = make([]Response, 0)
	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

import (
	"regexp"
	"strings"
	"sync/atomic"
)

type (
	Comment struct {
		Path string // /app/controllers/controller.go
		Line int    // 12
		Code string // type UserController struct{
		Name string // UserController
		Pkg  string // sketch/app/controllers

		Errors    []*Error    // Error codes.
		Headers   []*Header   // Request Header.
		Request   *Request    // Request params.
		Responses []*Response // Response list.

		offset int32
		text   []string
	}
)

func NewComment() *Comment {
	return (&Comment{
		text: make([]string, 0),
	}).init()
}

func (o *Comment) AddAnnotation(a Annotation) {
	switch a.GetType() {
	case AnnotationError:
		o.Errors = append(o.Errors, NewError(a))

	case AnnotationHeader:
		o.Headers = append(o.Headers, NewHeader(a))

	case AnnotationRequest:
		o.Request = NewRequest(a)

	case AnnotationResponse, AnnotationResponseData:
		o.Responses = append(o.Responses, NewResponse(a, ResponseData))

	case AnnotationResponseList:
		o.Responses = append(o.Responses, NewResponse(a, ResponseList))

	case AnnotationResponsePaging:
		o.Responses = append(o.Responses, NewResponse(a, ResponsePaging))
	}
}

func (o *Comment) AddText(str string) {
	n := atomic.AddInt32(&o.offset, 1)

	if n == 1 {
		str = strings.TrimPrefix(str, o.Name)
	}

	str = regexp.MustCompile(`[.]+$`).ReplaceAllString(str, "")

	if strings.TrimSpace(str) == "" {
		return
	}

	o.text = append(o.text, str)
}

func (o *Comment) GetDescription() (s string) {
	if len(o.text) > 1 {
		s = strings.TrimSpace(strings.Join(o.text[1:], "\n"))
	}
	return
}

func (o *Comment) GetTitle() string {
	if len(o.text) > 0 {
		return strings.TrimSpace(o.text[0])
	}
	return o.Name
}

func (o *Comment) init() *Comment {
	o.Errors = make([]*Error, 0)
	o.Headers = make([]*Header, 0)
	o.Responses = make([]*Response, 0)
	return o
}

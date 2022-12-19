// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package base

import (
	"fmt"
	"github.com/fuyibing/gdoc/conf"
	"net/http"
	"regexp"
	"strings"
)

type (
	Method interface {
		GetComment() *Comment
		GetContentType() ContentType
		GetController() Controller
		GetRequestMethod() string
		GetRequestUrl() string
		GetSortKey() string
	}

	method struct {
		comment    *Comment
		controller Controller

		ContentType               ContentType
		Name                      string
		RequestMethod, RequestUrl string
	}
)

func NewMethod(controller Controller, name string) Method {
	return (&method{
		comment:    NewComment(),
		controller: controller,

		ContentType:   ContentJson,
		RequestMethod: http.MethodGet,
	}).init(name)
}

func (o *method) GetComment() *Comment        { return o.comment }
func (o *method) GetContentType() ContentType { return o.ContentType }
func (o *method) GetController() Controller   { return o.controller }
func (o *method) GetRequestMethod() string    { return o.RequestMethod }
func (o *method) GetRequestUrl() string       { return o.getRequestUrl() }
func (o *method) GetSortKey() string          { return o.getSortKey() }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *method) getRequestUrl() (str string) {
	if str = fmt.Sprintf("%s%s", o.controller.GetPrefix(), o.RequestUrl); str == "" {
		str = "/"
		return
	}

	return
}

func (o *method) getSortKey() string {
	return fmt.Sprintf("%s_%s", o.getRequestUrl(), o.RequestMethod)
}

func (o *method) init(name string) *method {
	if m := conf.Regex.SourceRouteMethod.FindStringSubmatch(name); len(m) == 2 {
		o.RequestMethod = strings.ToUpper(m[1])
		return o
	}

	if m := conf.Regex.SourceRoute.FindStringSubmatch(name); len(m) == 3 {
		o.RequestMethod = strings.ToUpper(m[1])
		o.RequestUrl = m[2]

		// Change letter of upper to lower.
		o.RequestUrl = regexp.MustCompile(`[A-Z]`).ReplaceAllStringFunc(o.RequestUrl, func(s string) string {
			return fmt.Sprintf("/%s", strings.ToLower(s))
		})

		// Change underline as slash.
		o.RequestUrl = strings.ReplaceAll(o.RequestUrl, "_", "/")

		// Remove double slashes or prefix or suffix.
		o.RequestUrl = regexp.MustCompile(`/+`).ReplaceAllString(o.RequestUrl, "/")
		o.RequestUrl = strings.TrimPrefix(o.RequestUrl, "/")
		o.RequestUrl = strings.TrimSuffix(o.RequestUrl, "/")

		// Prefix with slashes.
		if o.RequestUrl != "" {
			o.RequestUrl = fmt.Sprintf("/%s", o.RequestUrl)
		}

		return o
	}

	return o
}

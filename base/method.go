// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package base

import (
	"fmt"
	"github.com/fuyibing/gdoc/config"
	"net/http"
	"regexp"
	"strings"
)

type (
	Method interface {
		GetComment() Comment
		GetController() Controller
		GetName() string
		GetRequestMethod() string
		GetRequestUrl() string
	}

	method struct {
		Comment    Comment
		Controller Controller

		Name          string // eg. PostUser
		RequestMethod string // eg. POST
		RequestUrl    string // eg. /user
	}
)

func NewMethod(controller Controller, name string) Method {
	return (&method{
		Controller: controller,
		Comment:    NewComment(name),
		Name:       name,
	}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *method) GetComment() Comment       { return o.Comment }
func (o *method) GetController() Controller { return o.Controller }
func (o *method) GetName() string           { return o.Name }
func (o *method) GetRequestMethod() string  { return o.RequestMethod }
func (o *method) GetRequestUrl() string     { return o.RequestUrl }

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *method) init() *method {
	// Method only.
	if config.Regex.GetRouteMethod().MatchString(o.Name) {
		o.RequestMethod = strings.ToUpper(o.Name)
		return o
	}

	// Method & URI.
	if m := config.Regex.GetRouteUrl().FindStringSubmatch(o.Name); len(m) == 3 {
		o.RequestMethod = strings.ToUpper(m[1])
		o.RequestUrl = m[2]

		// Letter from upper to lower with slashes
		// and change underline as slash.
		o.RequestUrl = regexp.MustCompile(`[A-Z]`).ReplaceAllStringFunc(regexp.MustCompile(`_+`).ReplaceAllString(o.RequestUrl, "/"), func(s string) string {
			return fmt.Sprintf("/%s", strings.ToLower(s))
		})

		o.RequestUrl = regexp.MustCompile(`/+`).ReplaceAllString(o.RequestUrl, "/")
		o.RequestUrl = strings.TrimPrefix(o.RequestUrl, "/")
		o.RequestUrl = strings.TrimSuffix(o.RequestUrl, "/")

		if o.RequestUrl != "" {
			o.RequestUrl = fmt.Sprintf("/%s", o.RequestUrl)
		}
		return o
	}

	// Default
	// request method.
	o.RequestMethod = http.MethodGet
	return o
}

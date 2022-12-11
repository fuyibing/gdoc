// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package base

import (
	"strings"
)

type (
	// Annotation
	// used in code comment.
	Annotation interface {
		GetFirst() string
		GetKey() string
		GetLine() int
		GetType() AnnotationType
		GetValue(i int) string
	}

	annotation struct {
		Key, Value string
		Line       int

		Length int
		Values []string
		Type   AnnotationType
	}
)

// NewAnnotation
//
//   NewAnnotation("Request", "app/logics/example.LoginRequest")
//   NewAnnotation("Response", "app/logics/example.LoginResponse")
func NewAnnotation(line int, key, value string) Annotation {
	return (&annotation{
		Key: key, Value: value, Line: line,
		Values: make([]string, 0),
	}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *annotation) GetFirst() string        { return o.getValue(0) }
func (o *annotation) GetKey() string          { return o.Key }
func (o *annotation) GetLine() int            { return o.Line }
func (o *annotation) GetType() AnnotationType { return o.Type }
func (o *annotation) GetValue(i int) string   { return o.getValue(i) }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *annotation) getValue(i int) string {
	if i < o.Length {
		return o.Values[i]
	}
	return ""
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *annotation) init() *annotation {
	o.Type = AnnotationType(o.Key)

	if o.Value != "" {
		for _, s := range strings.Split(o.Value, ",") {
			s = strings.TrimSpace(s)
			s = strings.TrimPrefix(s, `"`)
			s = strings.TrimSuffix(s, `"`)
			s = strings.TrimSpace(s)
			o.Values = append(o.Values, s)
		}
		o.Length = len(o.Values)
	}

	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

import (
	"strconv"
	"strings"
)

type (
	Annotation interface {
		GetBool(d bool) bool
		GetFirst() string
		GetName() string
		GetType() AnnotationType
		GetValue(i int) string
		GetValues(i int) string
	}

	annotation struct {
		Length      int
		Name, Param string
		Type        AnnotationType
		Values      []string
	}
)

func NewAnnotation(name, param string) Annotation {
	return (&annotation{
		Name:  name,
		Param: param,
		Type:  AnnotationType(name),
	}).init()
}

func (o *annotation) GetBool(d bool) bool {
	if s := o.getValue(0); s != "" {
		if b, be := strconv.ParseBool(s); be == nil {
			return b
		}
	}
	return d
}

func (o *annotation) GetFirst() string        { return o.getValue(0) }
func (o *annotation) GetName() string         { return o.Name }
func (o *annotation) GetType() AnnotationType { return o.Type }
func (o *annotation) GetValue(i int) string   { return o.getValue(i) }
func (o *annotation) GetValues(i int) string  { return o.getValues(i) }

func (o *annotation) init() *annotation {
	o.Values = make([]string, 0)

	if o.Param != "" {
		for _, s := range strings.Split(o.Param, ",") {
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

func (o *annotation) getValue(i int) (s string) {
	if i < o.Length {
		s = o.Values[i]
	}
	return
}

func (o *annotation) getValues(i int) (s string) {
	if i < o.Length {
		s = strings.Join(o.Values[i:], ", ")
	}
	return
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

import (
	"fmt"
	"sync"
)

type (
	Controller interface {
		GetComment() *Comment
		GetMapping() Mapping
		GetMethod(k string) Method
		GetMethods() map[string]Method
		GetPrefix() string
		GetSortKey() string
		SetMethod(k string, m Method)
		SetPrefix(s string)
	}

	controller struct {
		comment *Comment
		mapping Mapping
		methods map[string]Method
		mu      sync.RWMutex
		prefix  string
	}
)

func NewController(mapping Mapping) Controller {
	return (&controller{
		mapping: mapping,
		comment: NewComment(),
	}).init()
}

func (o *controller) GetComment() *Comment          { return o.comment }
func (o *controller) GetMapping() Mapping           { return o.mapping }
func (o *controller) GetMethod(k string) Method     { return o.getMethod(k) }
func (o *controller) GetMethods() map[string]Method { return o.methods }
func (o *controller) GetPrefix() string             { return o.prefix }
func (o *controller) GetSortKey() string            { return o.getSortKey() }
func (o *controller) SetMethod(k string, m Method)  { o.setMethod(k, m) }
func (o *controller) SetPrefix(s string)            { o.prefix = s }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *controller) getMethod(k string) Method {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if m, ok := o.methods[k]; ok {
		return m
	}
	return nil
}

func (o *controller) getSortKey() string {
	return fmt.Sprintf("%s/%s", o.prefix, o.comment.Name)
}

func (o *controller) init() *controller {
	o.mu = sync.RWMutex{}

	o.methods = make(map[string]Method)
	return o
}

func (o *controller) setMethod(k string, m Method) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.methods[k] = m
}

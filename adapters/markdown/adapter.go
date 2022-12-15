// author: wsfuyibing <websearch@163.com>
// date: 2022-12-15

package markdown

import (
	"github.com/fuyibing/gdoc/adapters/markdown/tpl"
	"github.com/fuyibing/gdoc/base"
)

type (
	// Adapter
	// use to create markdown file to specified directory.
	Adapter struct {
		mapping base.Mapping
	}
)

// New
// return adapter instance.
func New(mapping base.Mapping) *Adapter {
	return (&Adapter{mapping: mapping}).init()
}

// Run
// adapter instance.
func (o *Adapter) Run() {
	o.runReadme()
	o.runMethod()
}

func (o *Adapter) init() *Adapter {
	return o
}

func (o *Adapter) runMethod() {
	for _, c := range o.mapping.GetControllers() {
		for _, m := range c.GetMethods() {
			tpl.NewMethod(o.mapping, m).Save()
		}
	}
}

func (o *Adapter) runReadme() {
	tpl.NewReadme(o.mapping).Save()
}

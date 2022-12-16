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

// Initialize
// instance fields.
func (o *Adapter) init() *Adapter {
	return o
}

// Run method.
func (o *Adapter) runMethod() {
	for _, c := range o.mapping.GetControllers() {
		for _, m := range c.GetMethods() {
			// Ignore
			// method markdown.
			if m.GetComment().Ignored {
				continue
			}

			// Save method markdown.
			tpl.NewMethod(o.mapping, m).Save()
		}
	}
}

// Run readme.
func (o *Adapter) runReadme() {
	tpl.NewReadme(o.mapping).Save()
}

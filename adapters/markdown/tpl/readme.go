// author: wsfuyibing <websearch@163.com>
// date: 2022-12-15

package tpl

import (
	"fmt"
	"github.com/fuyibing/gdoc/adapters/markdown/i18n"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/conf"
	"sort"
	"strings"
)

type (
	Readme struct {
		mapping base.Mapping

		nl, sep   string
		templates []string
	}
)

// NewReadme
// create and return instance.
func NewReadme(mapping base.Mapping) *Readme {
	return (&Readme{
		mapping:   mapping,
		templates: make([]string, 0),
	}).init()
}

// Save
// markdown file to specified directory.
func (o *Readme) Save() {
	for _, call := range []func(){
		o.HeaderTitle,
		o.Info,
		o.Description,

		o.Menu,

		o.save,
	} {
		call()
	}
}

// /////////////////////////////////////////////////////////////
// Internal methods
// /////////////////////////////////////////////////////////////

func (o *Readme) Description() {
	if str := conf.Config.Description; str != "" {
		list := make([]string, 0)
		for _, s := range strings.Split(str, o.nl) {
			if s = strings.TrimSpace(s); s != "" {
				list = append(list, fmt.Sprintf("> %s", s))
			}
		}
		o.templates = append(o.templates, strings.Join(list, o.nl))
	}
}

func (o *Readme) HeaderTitle() {
	o.templates = append(o.templates, fmt.Sprintf("# %v", conf.Config.Name))
}

func (o *Readme) Info() {
	o.templates = append(o.templates, strings.Join([]string{
		fmt.Sprintf("**%v** : `%s`", i18n.Lang("Host"), conf.Config.Deploy.Full()),
		fmt.Sprintf("**%v** : `%v`", i18n.Lang("Updated"), o.mapping.GetLastUpdated()),
	}, "<br />"))
}

func (o *Readme) Menu() {
	o.templates = append(o.templates, fmt.Sprintf("### %s",
		i18n.Lang("Table Contents"),
	))
	o.controllers()
}

// /////////////////////////////////////////////////////////////
// Access methods.
// /////////////////////////////////////////////////////////////

func (o *Readme) controller(c base.Controller) []string {
	var (
		list = make([]string, 0)
		item = o.methods(c)
	)

	if n := len(item); n > 0 {
		list = append(list,
			fmt.Sprintf("* *%s* <small>(%d)</small>", c.GetComment().GetTitle(), n),
			strings.Join(item, o.nl),
		)
	}

	return list
}

func (o *Readme) controllers() {
	var (
		keys = make([]string, 0)
		maps = make(map[string]base.Controller)
		list = make([]string, 0)
	)

	// Build
	// controllers key.
	for _, c := range o.mapping.GetControllers() {
		k := c.GetSortKey()
		keys = append(keys, k)
		maps[k] = c
	}

	// Sort controllers key.
	sort.Strings(keys)

	// Iterate controller by sorted.
	for _, k := range keys {
		if ls := o.controller(maps[k]); len(ls) > 0 {
			list = append(list, ls...)
		}
	}

	// Append to template.
	if len(list) > 0 {
		o.templates = append(o.templates, strings.Join(list, o.nl))
	}
}

func (o *Readme) methods(c base.Controller) []string {
	var (
		keys = make([]string, 0)
		maps = make(map[string]base.Method)
		list = make([]string, 0)
	)

	// Build
	// method key.
	for _, m := range c.GetMethods() {
		// Skip ignored method.
		if m.GetComment().Ignored {
			continue
		}

		// Add sortable key.
		k := m.GetSortKey()
		keys = append(keys, k)
		maps[k] = m
	}

	// Sort controllers key.
	sort.Strings(keys)

	// Iterate controller by sorted.
	for _, k := range keys {
		m := maps[k]
		l := linkName(m)
		list = append(list, fmt.Sprintf(
			"  * [%s](.%s) - <small>`%s`</small> <small>`%s`</small>",
			m.GetComment().GetTitle(), l,
			m.GetRequestMethod(),
			m.GetRequestUrl(),
		))
	}

	return list
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *Readme) init() *Readme {
	o.nl = "\n"
	o.sep = "\n\n"
	return o
}

func (o *Readme) save() {
	// Append source info
	// to markdown files end.
	o.templates = append(o.templates, "----",
		strings.Join([]string{
			fmt.Sprintf("* `/%v` - match module name", "go.md"),
			fmt.Sprintf("* `/%v` - match application configurations", conf.Path.GetDocumentJsonFile()),
			fmt.Sprintf("* `%v` - controller files location", conf.Path.GetControllerPath()),
			fmt.Sprintf("* `%v` - document storage location", conf.Path.GetDocumentPath()),
		}, o.nl),
	)

	// Call save handler
	// to specified file.
	conf.Path.SavePath(
		fmt.Sprintf("%s%s/README.md",
			conf.Path.GetBasePath(),
			conf.Path.GetDocumentPath(),
		), strings.Join(o.templates,
			o.sep,
		),
	)
}

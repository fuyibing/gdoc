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
		mapping   base.Mapping
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
		o.header,
		o.description,
		o.deploy,
		o.menu,
		o.save,
	} {
		call()
	}
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Readme) deploy() {
	o.templates = append(o.templates, "----", strings.Join([]string{
		fmt.Sprintf("* %v: `%v`", i18n.Lang("Deploy Port"), conf.Config.Deploy.Port),
		fmt.Sprintf("* %v: `%v`", i18n.Lang("Deploy Host"), conf.Config.Deploy.Host),
		fmt.Sprintf("* %v: `%v`", i18n.Lang("Updated"), o.mapping.GetLastUpdated()),
	}, "\n"))
}

func (o *Readme) description() {
	if str := conf.Config.Description; str != "" {
		o.templates = append(o.templates, "----", str)
	}
}

func (o *Readme) header() {
	o.templates = append(o.templates,
		fmt.Sprintf("# %s", conf.Config.Name),
	)
}

func (o *Readme) init() *Readme {
	return o
}

func (o *Readme) menu() {
	var (
		list = make([]string, 0)
		ck   = make([]string, 0)
		cm   = make(map[string]base.Controller)
	)

	for _, c := range o.mapping.GetControllers() {
		k := c.GetSortKey()

		cm[k] = c
		ck = append(ck, c.GetSortKey())
	}

	sort.Strings(ck)
	for _, k := range ck {
		c := cm[k]
		list = append(list,
			fmt.Sprintf("* %s `%s`", c.GetComment().GetTitle(), c.GetComment().Name),
		)

		func(c base.Controller) {
			ms := make([]string, 0)
			mm := make(map[string]base.Method)

			for _, m := range c.GetMethods() {
				if m.GetComment().Ignored {
					continue
				}
				mk := m.GetSortKey()
				mm[mk] = m
				ms = append(ms, mk)
			}

			sort.Strings(ms)
			for _, mk := range ms {
				m := mm[mk]
				l := linkName(m)
				list = append(list,
					fmt.Sprintf("  * [%s](.%s) `%s` `%s`",
						m.GetComment().GetTitle(),
						l,
						m.GetRequestMethod(),
						m.GetRequestUrl(),
					),
				)
			}
		}(c)
	}

	o.templates = append(o.templates,
		fmt.Sprintf("### %s", i18n.Lang("Controller Menu")),
		strings.Join(list, "\n"),
	)
}

func (o *Readme) save() {
	conf.Path.SavePath(
		fmt.Sprintf("%s%s/README.md",
			conf.Path.GetBasePath(),
			conf.Path.GetDocumentPath(),
		),
		strings.Join(o.templates, "\n\n"),
	)
}

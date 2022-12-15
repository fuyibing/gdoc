// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package reflectors

import (
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/conf"
	"os/exec"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

type (
	Reflection interface {
		Configure() Reflection
		Make() Reflection
	}

	mapper struct {
		alias string
		name  string
	}

	reflection struct {
		mapping base.Mapping

		packageIndex  int32
		packageList   []string
		packageMapper map[string]string
		structMapper  map[string]*mapper
	}
)

func New(mapping base.Mapping) Reflection {
	return (&reflection{mapping: mapping}).init()
}

func (o *reflection) Configure() Reflection {
	// Iterate controllers.
	for _, c := range o.mapping.GetControllers() {
		// Iterate methods.
		for _, m := range c.GetMethods() {
			// Request struct.
			if r := m.GetComment().Request; r != nil {
				o.mapper(r.Key, r.Pkg, r.Name)
			}
			// Response struct list.
			for _, r := range m.GetComment().Responses {
				o.mapper(r.Key, r.Pkg, r.Name)
			}
		}
	}

	// Create tmp files.
	o.save()
	return o
}

func (o *reflection) Make() Reflection {
	cmd := exec.Command("go", "run", fmt.Sprintf(
		"%s%s%s/main.go",
		conf.Path.GetBasePath(),
		conf.Path.GetDocumentPath(),
		conf.Path.GetTmpPath(),
	))

	buf, err := cmd.Output()

	if err != nil {
		conf.Debugger.Error("make error: %v", err)
	} else {
		conf.Debugger.Info("make done: %s", buf)
	}

	return o
}

func (o *reflection) init() *reflection {
	o.packageList = make([]string, 0)
	o.packageMapper = make(map[string]string)
	o.structMapper = make(map[string]*mapper)
	return o
}

func (o *reflection) mapper(key, pkg, name string) {
	var (
		alias string
		ok    bool
	)

	// Generate alias.
	if alias, ok = o.packageMapper[pkg]; !ok {
		n := atomic.AddInt32(&o.packageIndex, 1)
		alias = fmt.Sprintf("a%d", n)

		o.packageList = append(o.packageList, pkg)
		o.packageMapper[pkg] = alias
	}

	// Update struct mapper.
	o.structMapper[key] = &mapper{
		alias: alias, name: name,
	}
}

func (o *reflection) save() {
	var (
		path = fmt.Sprintf("%s%s%s/main.go",
			conf.Path.GetBasePath(),
			conf.Path.GetDocumentPath(),
			conf.Path.GetTmpPath(),
		)
		text = []string{
			fmt.Sprintf("// Reflection template."),
			fmt.Sprintf("// Do not edit this file: %s", time.Now().Format("2006-01-02 15:04:05")),
			"",
			"package main",
		}
	)

	// Imports.
	sort.Strings(o.packageList)
	text = append(text,
		"",
		"import (",
		"    \"github.com/fuyibing/gdoc/reflectors\"",
		"",
	)
	for _, k := range o.packageList {
		text = append(text,
			fmt.Sprintf("    %s \"%s\"", o.packageMapper[k], k),
		)
	}
	text = append(text, ")")

	// Main function.
	text = append(text, "",
		"func main() {",
		"    ref := reflectors.Parser()",
		fmt.Sprintf("    ref.BasePath = \"%s\"", conf.Path.GetBasePath()),
		fmt.Sprintf("    ref.ControllerPath = \"%s\"", conf.Path.GetControllerPath()),
		fmt.Sprintf("    ref.DocumentPath = \"%s\"", conf.Path.GetDocumentPath()),
		fmt.Sprintf("    ref.DocumentJsonFile = \"%s\"", conf.Path.GetDocumentJsonFile()),
		fmt.Sprintf("    ref.TmpPath = \"%s\"", conf.Path.GetTmpPath()),
		"",
		"    for k, p := range map[string]interface{}{",
	)

	for k, v := range o.structMapper {
		text = append(text,
			fmt.Sprintf("        \"%s\": &%s.%s{},", k, v.alias, v.name),
		)
	}

	text = append(text,
		"    }{",
		"        ref.Parse(k, p)",
		"    }",
		"",
		"    ref.Save()",
		"}",
	)

	// Save file.
	conf.Path.SavePath(path, strings.Join(text, "\n"))
}
// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package managers

import (
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/config"
	"os"
)

type (
	Scanner interface {
		GetController(key string) base.Controller
		GetControllers() map[string]base.Controller
		Scan() error
		SetController(key string, controller base.Controller)
	}

	scanner struct {
		Manager          Management
		ControllerList   []string
		ControllerMapper map[string]base.Controller
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *scanner) GetController(key string) base.Controller   { return o.getController(key) }
func (o *scanner) GetControllers() map[string]base.Controller { return o.ControllerMapper }
func (o *scanner) SetController(k string, c base.Controller)  { o.setController(k, c) }
func (o *scanner) Scan() error                                { return o.scan("") }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *scanner) getController(key string) base.Controller {
	if c, ok := o.ControllerMapper[key]; ok {
		return c
	}
	return nil
}

func (o *scanner) scan(prefix string) (err error) {
	var (
		ds   []os.DirEntry
		path = fmt.Sprintf("%s%s", config.Path.GetBase(), config.Path.GetController())
	)

	if prefix != "" {
		path = fmt.Sprintf("%s%s", path, prefix)
	}

	o.Manager.GetLogger().Info("[scan] read directory at %s", path)

	if ds, err = os.ReadDir(path); err != nil {
		return
	}

	for _, d := range ds {
		if config.Regex.GetHiddenFile().MatchString(d.Name()) {
			continue
		}

		if d.IsDir() {
			if err = o.scan(fmt.Sprintf("%s/%s", prefix, d.Name())); err != nil {
				return
			}
			continue
		}

		if config.Regex.GetSourceFile().MatchString(d.Name()) {
			if err = NewReader(o.Manager, prefix, d.Name()).Run(); err != nil {
				return err
			}
		}
	}

	return
}

func (o *scanner) setController(key string, controller base.Controller) {
	o.ControllerList = append(o.ControllerList, fmt.Sprintf("%s.%s", controller.GetPackage(), controller.GetName()))
	o.ControllerMapper[key] = controller
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *scanner) init() *scanner {
	o.ControllerList = make([]string, 0)
	o.ControllerMapper = make(map[string]base.Controller)
	return o
}

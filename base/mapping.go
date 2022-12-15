// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

import (
	"sync"
)

var (
	Mapper Mapping
)

type (
	Mapping interface {
		GetController(k string) Controller
		GetControllers() map[string]Controller
		SetController(k string, c Controller)
	}

	mapping struct {
		controllers map[string]Controller
		mu          sync.RWMutex
	}
)

func (o *mapping) GetController(k string) Controller     { return o.getController(k) }
func (o *mapping) GetControllers() map[string]Controller { return o.controllers }
func (o *mapping) SetController(k string, c Controller)  { o.setController(k, c) }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *mapping) getController(k string) Controller {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if c, ok := o.controllers[k]; ok {
		return c
	}
	return nil
}

func (o *mapping) init() *mapping {
	o.controllers = make(map[string]Controller)
	o.mu = sync.RWMutex{}
	return o
}

func (o *mapping) setController(k string, c Controller) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.controllers[k] = c
}

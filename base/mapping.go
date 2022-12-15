// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

import (
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gdoc/conf"
	"os"
	"sync"
	"time"
)

var (
	Mapper Mapping
)

type (
	Mapping interface {
		GetController(k string) Controller
		GetControllers() map[string]Controller
		GetLastUpdated() string
		LoadTmpCode(k string) []byte
		LoadTmpItem(k string) []*Item
		SetController(k string, c Controller)
	}

	mapping struct {
		controllers map[string]Controller
		mu          sync.RWMutex
		updated     string
	}
)

func (o *mapping) GetController(k string) Controller     { return o.getController(k) }
func (o *mapping) GetControllers() map[string]Controller { return o.controllers }
func (o *mapping) GetLastUpdated() string                { return o.updated }
func (o *mapping) LoadTmpCode(k string) []byte           { return o.loadTmpCode(k) }
func (o *mapping) LoadTmpItem(k string) []*Item          { return o.loadTmpItem(k) }
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
	o.updated = time.Now().Format("2006-01-01 15:04")
	return o
}

func (o *mapping) loadTmp(name string) []byte {
	if buf, err := os.ReadFile(fmt.Sprintf("%s%s%s/%s",
		conf.Path.GetBasePath(),
		conf.Path.GetDocumentPath(),
		conf.Path.GetTmpPath(),
		name,
	)); err == nil {
		return buf
	}
	return nil
}

func (o *mapping) loadTmpCode(k string) []byte {
	return o.loadTmp(conf.Path.GenerateCodeFile(k))
}

func (o *mapping) loadTmpItem(k string) (items []*Item) {
	if buf := o.loadTmp(conf.Path.GenerateItemFile(k)); buf != nil {
		items = make([]*Item, 0)
		_ = json.Unmarshal(buf, &items)
	}
	return
}

func (o *mapping) setController(k string, c Controller) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.controllers[k] = c
}

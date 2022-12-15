// author: wsfuyibing <websearch@163.com>
// date: 2022-12-15

package reflectors

import (
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gdoc/conf"
	"reflect"
	"sync"
)

type (
	ParserManager struct {
		BasePath         string
		ControllerPath   string
		DocumentPath     string
		DocumentJsonFile string
		TmpPath          string

		Structs map[string]*Struct
		mu      sync.RWMutex
	}
)

// Parser
//
// Create instance in target application, and called in
// terminal by coder.
func Parser() *ParserManager {
	return (&ParserManager{}).init()
}

// Parse
// call reflect manager and convert fields into tmp files.
func (o *ParserManager) Parse(key string, ptr interface{}) {
	o.mu.Lock()

	// Return
	// if parse already.
	if _, ok := o.Structs[key]; ok {
		o.mu.Unlock()
		return
	}

	// Add to memory.
	s := NewStruct(o)
	o.Structs[key] = s
	o.mu.Unlock()

	// Iterate progress.
	x := reflect.ValueOf(ptr).Elem()
	s.Iterate(x)
}

// Save
// as json tmp files.
func (o *ParserManager) Save() {
	for k, v := range o.Structs {
		func(fk string, fp *Struct) {
			for mk, mp := range map[string]interface{}{
				conf.Path.GenerateCodeFile(fk): fp.Map(),
				conf.Path.GenerateItemFile(fk): fp.Items(),
			} {
				o.save(mk, mp)
			}
		}(k, v)
	}
}

func (o *ParserManager) init() *ParserManager {
	conf.Path.SetBasePath(o.BasePath)
	conf.Path.SetControllerPath(o.ControllerPath)
	conf.Path.SetDocumentPath(o.DocumentPath)
	conf.Path.SetDocumentJsonFile(o.DocumentJsonFile)
	conf.Path.SetTmpPath(o.TmpPath)
	conf.Config.Load()

	o.Structs = make(map[string]*Struct)
	o.mu = sync.RWMutex{}
	return o
}

func (o *ParserManager) save(name string, ptr interface{}) {
	var (
		path     = fmt.Sprintf("%s%s%s/%s", o.BasePath, o.DocumentPath, o.TmpPath, name)
		buf, err = json.MarshalIndent(ptr, "", "    ")
	)

	if err != nil {
		conf.Debugger.Error("[parse] %v", err)
		return
	}

	conf.Path.SavePath(path, string(buf))
}

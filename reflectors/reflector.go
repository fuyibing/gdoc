// author: wsfuyibing <websearch@163.com>
// date: 2022-12-05

package reflectors

import (
	"fmt"
	"reflect"
)

type (
	Reflection interface {
		Blocks() map[string]Block
		Parse(key string, ptr interface{}) error
	}

	reflection struct {
		blocks map[string]Block
	}
)

func New() Reflection {
	return (&reflection{}).init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *reflection) Blocks() map[string]Block                { return o.blocks }
func (o *reflection) Parse(key string, ptr interface{}) error { return o.parse(key, ptr) }

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *reflection) init() *reflection {
	o.blocks = make(map[string]Block)
	return o
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *reflection) parse(key string, ptr interface{}) error {
	if _, ok := o.blocks[key]; ok {
		return nil
	}

	r := reflect.ValueOf(ptr)
	if r.Kind() != reflect.Ptr {
		return fmt.Errorf("type of instance for reflector must be pointer")
	}

	x := NewBlock(o)
	o.blocks[key] = x

	return x.Parse(r.Elem())
}

func (o *reflection) saveJson(key string, block Block) error {
	return nil
}

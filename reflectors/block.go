// author: wsfuyibing <websearch@163.com>
// date: 2022-12-06

package reflectors

import (
	"github.com/fuyibing/gdoc/base"
	"reflect"
	"sort"
)

type (
	Block interface {
		Export() map[string]interface{}
		Fields() map[string]Field
		Parse(v reflect.Value) error
		Reflection() Reflection
		ToList() []*base.Item
	}

	block struct {
		reflection  Reflection
		fieldList   []string
		fieldMapper map[string]Field
	}
)

func NewBlock(reflection Reflection) Block {
	return (&block{reflection: reflection}).
		init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *block) Export() map[string]interface{} { return o.export() }
func (o *block) Fields() map[string]Field       { return o.fieldMapper }
func (o *block) Parse(v reflect.Value) error    { return o.parse(v) }
func (o *block) Reflection() Reflection         { return o.reflection }
func (o *block) ToList() []*base.Item           { return o.toList() }

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *block) init() *block {
	o.fieldList = make([]string, 0)
	o.fieldMapper = make(map[string]Field)
	return o
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *block) export() (res map[string]interface{}) {
	res = make(map[string]interface{})
	sort.Strings(o.fieldList)
	for _, k := range o.fieldList {
		if f, ok := o.fieldMapper[k]; ok {
			res[f.GetSortKey()] = f.Export()
		}
	}
	return res
}

func (o *block) parse(v reflect.Value) (err error) {
	for i := 0; i < v.NumField(); i++ {
		iv := v.Field(i)
		isf := v.Type().Field(i)

		if isf.Anonymous {
			if iv.Kind() == reflect.Struct {
				if err = o.parse(iv); err != nil {
					return
				}
				continue
			}

			if iv.Kind() == reflect.Ptr {
				err = ErrAnonymousPointerNotAllowed
				return
			}

			err = ErrAnonymousUnknown
			return
		}

		if err = o.parseField(iv, isf); err != nil {
			return
		}
	}
	return
}

func (o *block) parseField(v reflect.Value, sf reflect.StructField) (err error) {
	f := NewField(o)

	if err = f.Parse(v, sf); err != nil {
		return
	}

	if f.GetKey() != "" {
		k := f.GetSortKey()
		o.fieldList = append(o.fieldList, k)
		o.fieldMapper[k] = f
	}

	return
}

func (o *block) toList() (list []*base.Item) {
	sort.Strings(o.fieldList)
	list = make([]*base.Item, 0)
	for _, k := range o.fieldList {
		if f, ok := o.fieldMapper[k]; ok {
			list = append(list, f.Item())
		}
	}
	return
}

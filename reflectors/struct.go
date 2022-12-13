// author: wsfuyibing <websearch@163.com>
// date: 2022-12-12

package reflectors

import (
	"github.com/fuyibing/gdoc/base"
	"reflect"
	"regexp"
	"sort"
	"sync"
)

type (
	// Struct
	// 结构体.
	Struct struct {
		FieldList   []string
		FieldMapper map[string]*Field

		mu         sync.RWMutex
		reflection *Reflection
	}
)

func NewStruct(reflection *Reflection) *Struct {
	return &Struct{
		FieldList:   make([]string, 0),
		FieldMapper: make(map[string]*Field),

		mu:         sync.RWMutex{},
		reflection: reflection,
	}
}

// Items
// 导出元素列表.
func (o *Struct) Items() []*base.Item {
	sort.Strings(o.FieldList)

	items := make([]*base.Item, 0)
	for _, k := range o.FieldList {
		if f, ok := o.FieldMapper[k]; ok {
			items = append(items, f.Item())
		}
	}
	return items
}

// Iterate
// 迭代字段.

func (o *Struct) Iterate(v reflect.Value) {
	o.reflection.Info("find struct: %v.%v", v.Type().PkgPath(), v.Type().Name())

	for i := 0; i < v.NumField(); i++ {
		iv := v.Field(i)
		isf := v.Type().Field(i)

		// 非匿名.
		if !isf.Anonymous {
			if regexp.MustCompile(`^[A-Z]`).MatchString(isf.Name) {
				o.field(iv, isf)
			}
			continue
		}

		// 匿名指针.
		if iv.Kind() == reflect.Ptr {
			o.Iterate(reflect.New(iv.Type().Elem()).Elem())
			continue
		}

		// 匿名结构体.
		if iv.Kind() == reflect.Struct {
			o.Iterate(v)
			continue
		}
	}
	return
}

// Map
// 导出键值对.
func (o *Struct) Map() map[string]interface{} {
	sort.Strings(o.FieldList)

	res := make(map[string]interface{})
	for _, k := range o.FieldList {
		if f, ok := o.FieldMapper[k]; ok {
			if f.Ignored {
				continue
			}
			res[f.Key] = f.Map()
		}
	}
	return res
}

// 解析字段.
func (o *Struct) field(v reflect.Value, sf reflect.StructField) {
	f := NewField(o, sf)
	f.Parse(v)

	o.mu.Lock()
	o.FieldList = append(o.FieldList, f.SortKey())
	o.FieldMapper[f.SortKey()] = f
	o.mu.Unlock()
}

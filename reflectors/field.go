// author: wsfuyibing <websearch@163.com>
// date: 2020-12-12

package reflectors

import (
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"reflect"
	"regexp"
	"strconv"
)

type (
	Field struct {
		s     *Struct
		child *Struct
		mock  string

		Array              bool
		Condition          string
		Ignored            bool
		Key, Name          string
		Kind               FieldKind
		Label, Description string
		Required           bool
		Type               string
		Value              interface{}
	}
)

func NewField(s *Struct, sf reflect.StructField) *Field {
	return (&Field{s: s}).init(&sf)
}

func (o *Field) Item() *base.Item {
	item := &base.Item{
		Array:       o.Array,
		Condition:   o.Condition,
		Description: o.Description,
		Ignored:     o.Ignored,
		Key:         o.Key,
		Kind:        int(o.Kind),
		Label:       o.Label,
		Name:        o.Name,
		Required:    o.Required,
		Type:        o.Type,
		Value:       o.Value,
	}
	if o.child != nil {
		item.Children = o.child.Items()
	}
	return item
}

func (o *Field) Map() interface{} {
	if o.Array {
		if o.child != nil {
			return []interface{}{
				o.child.Map(),
			}
		}
		return []interface{}{
			o.Value,
		}
	}

	if o.child != nil {
		return o.child.Map()
	}

	return o.Value
}

func (o *Field) Parse(v reflect.Value) {
	// Recursion.
	if v.Kind() == reflect.Struct {
		o.Type = v.Type().PkgPath() + "." + v.Type().Name()

		o.child = NewStruct(o.s.parser)
		o.child.Iterate(v)
		return
	}

	// Type redirect for pointer.
	if v.Kind() == reflect.Ptr {
		o.Parse(reflect.New(v.Type().Elem()).Elem())
		return
	}

	// Type redirect for slice.
	if v.Kind() == reflect.Slice {
		o.Array = true
		o.Parse(reflect.New(v.Type().Elem()).Elem())
		return
	}

	// Type explain on map.
	if v.Kind() == reflect.Map {
		o.Type = "MAP"
		o.Value = make(map[int]int)
		return
	}

	// System type.
	//
	// - reflect.Uintptr
	// - reflect.Complex64
	// - reflect.Complex128
	// - reflect.Array
	// - reflect.Chan
	// - reflect.Func
	// - reflect.UnsafePointer
	switch v.Kind() {
	case reflect.Interface:
		o.Value = "*"

	case reflect.Bool:
		o.Value = false

	case reflect.Float32, reflect.Float64:
		o.Value = 0

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		o.Value = 0

	case reflect.String:
		o.Value = ""
	}

	o.Type = v.Kind().String()
	o.parseMock()
	return
}

func (o *Field) SortKey() string {
	if o.Required {
		return fmt.Sprintf("0_%s", o.Key)
	}
	return fmt.Sprintf("1_%s", o.Key)
}

func (o *Field) init(sf *reflect.StructField) *Field {
	o.Name = sf.Name
	o.Key = sf.Name
	o.Kind = FieldJson

	for kind, tag := range map[FieldKind]string{
		FieldJson: "json",
		FieldForm: "form",
		FieldUrl:  "url",
	} {
		if s := sf.Tag.Get(tag); s != "" {
			o.Kind = kind
			o.Key = s
		}
	}

	o.initDesc(sf)
	o.initLabel(sf)
	o.initMock(sf)
	o.initIgnore()
	o.initValidate(sf)
	return o
}

func (o *Field) initDesc(sf *reflect.StructField) {
	if s := sf.Tag.Get("desc"); s != "" {
		o.Description = s
	}
}

func (o *Field) initIgnore() {
	if o.Key == "-" {
		o.Ignored = true
	}
}

func (o *Field) initLabel(sf *reflect.StructField) {
	o.Label = o.Name
	if s := sf.Tag.Get("label"); s != "" {
		o.Label = s
	}
}

func (o *Field) initMock(sf *reflect.StructField) {
	if s := sf.Tag.Get("mock"); s != "" {
		o.mock = s
	}
}

func (o *Field) initValidate(sf *reflect.StructField) {
	if s := sf.Tag.Get("validate"); s != "" {
		if regexp.MustCompile(`required`).MatchString(s) {
			o.Required = true

			s = regexp.MustCompile(`,\s*required`).ReplaceAllString(s, "")
			s = regexp.MustCompile(`required\s*,\s*`).ReplaceAllString(s, "")
		}
		o.Condition = s
	}
}

func (o *Field) parseMock() {
	if o.mock == "" {
		return
	}

	if o.Type == "bool" {
		if n, err := strconv.ParseBool(o.mock); err == nil {
			o.Value = n
		}
	}

	if regexp.MustCompile(`^float`).MatchString(o.Type) {
		if n, err := strconv.ParseFloat(o.mock, 64); err == nil {
			o.Value = n
		}
	}

	if regexp.MustCompile(`^int`).MatchString(o.Type) {
		if n, err := strconv.ParseInt(o.mock, 0, 64); err == nil {
			o.Value = n
		}
	}

	if regexp.MustCompile(`^uint`).MatchString(o.Type) {
		if n, err := strconv.ParseUint(o.mock, 0, 64); err == nil {
			o.Value = n
		}
	}

	if o.Type == "string" {
		o.Value = o.mock
		return
	}
}

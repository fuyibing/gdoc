// author: wsfuyibing <websearch@163.com>
// date: 2022-12-06

package reflectors

import (
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"reflect"
	"strconv"
)

type (
	// Field
	// 字段接口.
	Field interface {
		Block() Block
		Export() interface{}
		GetDescription() string
		GetKey() string
		GetLabel() string
		GetMock() string
		GetSortKey() string
		GetTypeName() string
		GetValidation() string
		IsArray() bool
		IsExec() bool
		IsRequired() bool
		IsSystemType() bool
		Item() *base.Item
		Parse(v reflect.Value, sf reflect.StructField) error
	}

	// 结构块.
	field struct {
		block Block

		isArray      bool // 是否为数组
		isExec       bool // 是否计算属性
		isRequired   bool // 是否必须
		isSystemType bool // 是否为系统类型

		child        Block       // 子类型.
		description  string      // 字段描述
		defaultValue interface{} // 默认值.
		key          string      // 字段键名, 如: user
		label        string      // 标签名称, 如: 邮箱地址
		mock         string      // 防造数据, 如: websearch@163.com
		typeName     string      // 类型名称, 如: int
		validation   string      // 校验选项
	}
)

func NewField(block Block) Field {
	return (&field{block: block}).
		init()
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *field) Block() Block                                        { return o.block }
func (o *field) Export() interface{}                                 { return o.export() }
func (o *field) GetDescription() string                              { return o.description }
func (o *field) GetKey() string                                      { return o.key }
func (o *field) GetLabel() string                                    { return o.label }
func (o *field) GetMock() string                                     { return o.mock }
func (o *field) GetSortKey() string                                  { return o.sortKey() }
func (o *field) GetTypeName() string                                 { return o.typeName }
func (o *field) GetValidation() string                               { return o.validation }
func (o *field) IsArray() bool                                       { return o.isArray }
func (o *field) IsExec() bool                                        { return o.isExec }
func (o *field) IsRequired() bool                                    { return o.isRequired }
func (o *field) IsSystemType() bool                                  { return o.isSystemType }
func (o *field) Item() *base.Item                                    { return o.item() }
func (o *field) Parse(v reflect.Value, sf reflect.StructField) error { return o.parse(v, sf) }

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *field) init() *field {
	return o
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *field) export() interface{} {
	// 1. 系统类型.
	//    当值类型为系统类型时, 此时为最细颗粒度.
	if o.isSystemType {
		// 1.1 数组.
		if o.isArray {
			return []interface{}{
				o.defaultValue,
			}
		}

		// 1.2 默认.
		return o.defaultValue
	}

	// 2. 用户类型.
	//    当指定类型非系统时(结构体).
	if o.isArray {
		return []interface{}{
			o.child.Export(),
		}
	}

	// 3. 子类型.
	return o.child.Export()
}

func (o *field) item() *base.Item {
	it := &base.Item{
		// Key:         o.key,
		// Array:       o.isArray,
		// Type:        o.typeName,
		// Required:    o.isRequired,
		// Validation:  o.validation,
		// Description: o.description,
		// Label:       o.label,
		// Mock:        o.mock,
	}
	if o.child != nil {
		it.Children = o.child.ToList()
	}
	return it
}

func (o *field) parse(v reflect.Value, sf reflect.StructField) error {
	o.label = sf.Name
	o.key = sf.Name

	o.parseKey(sf)
	o.parseTag(sf)
	o.parseValidate(sf)

	return o.parseType(v)
}

func (o *field) parseKey(sf reflect.StructField) {
	if s := sf.Tag.Get(TagJson); s != "" {
		if s == TagIgnored {
			s = ""
		}
		o.key = s
	}
}

func (o *field) parseTag(sf reflect.StructField) {
	// 1. 标签名称.
	if s := sf.Tag.Get(TagLabel); s != "" {
		o.label = s
	}

	// 2. 字段描述.
	if s := sf.Tag.Get(TagDescription); s != "" {
		o.description = s
	}

	// 3. 计算属性.
	if s := sf.Tag.Get(TagExec); s != "" {
		o.isExec, _ = strconv.ParseBool(s)
	}

	// 4. 模拟数据.
	if s := sf.Tag.Get(TagMock); s != "" {
		o.mock = s
	}
}

func (o *field) parseType(v reflect.Value) (err error) {
	o.isSystemType = true

	switch v.Kind() {
	case reflect.Struct:
		{
			o.child = NewBlock(o.block.Reflection())
			o.isSystemType = false
			o.typeName = "object"

			if err = o.child.Parse(v); err != nil {
				return err
			}
		}

	case reflect.Ptr, reflect.Slice:
		{
			if v.Kind() == reflect.Slice {
				o.isArray = true
			}
			if err = o.parseType(reflect.New(v.Type().Elem()).Elem()); err != nil {
				return err
			}
		}

	case reflect.Bool:
		{
			o.defaultValue = DefaultBoolValue
			o.typeName = "bool"
		}

	case reflect.Float32, reflect.Float64:
		{
			o.defaultValue = DefaultFloatValue
			o.typeName = "float"
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			o.defaultValue = DefaultIntValue
			o.typeName = "int"
		}

	case reflect.String:
		{
			o.defaultValue = DefaultStringValue
			o.typeName = "string"
		}

	case reflect.Interface:
		{
			o.defaultValue = DefaultInterfaceValue
			o.typeName = "*"
		}

	case reflect.Map:
		{
			o.defaultValue = DefaultMapValue
			o.typeName = "map"
		}

	default:
		err = ErrUnknownType
	}
	return
}

func (o *field) parseValidate(sf reflect.StructField) {
	s := sf.Tag.Get(TagValidate)

	if s == "" {
		return
	}

	if RegexpFieldValidate.MatchString(s) {
		o.isRequired = true
		o.validation = RegexpFieldValidate.ReplaceAllString(s, "")
	}
}

func (o *field) sortKey() string {
	if o.isRequired {
		return fmt.Sprintf("0_%s", o.key)
	}
	return fmt.Sprintf("1_%s", o.key)
}

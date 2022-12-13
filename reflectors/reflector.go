// author: wsfuyibing <websearch@163.com>
// date: 2022-12-05

// Package reflectors
//
// 通过结构体反射, 解析指定结构下的字段, 以及嵌套的结构体.
package reflectors

import (
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"os"
	"reflect"
	"regexp"
	"sync"
)

type (
	// Reflection
	// 反射入口.
	Reflection struct {
		Structs                        map[string]*Struct
		BasePath, StoragePath, TmpPath string

		directories map[string]bool
		mu          sync.RWMutex
	}
)

// New
// 创建反射.
func New() *Reflection {
	return &Reflection{
		Structs: make(map[string]*Struct),

		directories: make(map[string]bool),
		mu:          sync.RWMutex{},
	}
}

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

// Parse
// 解析结构体.
//
//   ref.Parse("pkg.Example", &pkg.Example{})
func (o *Reflection) Parse(key string, ptr interface{}) error {
	o.mu.Lock()

	if _, ok := o.Structs[key]; ok {
		o.mu.Unlock()
		return nil
	}

	s := NewStruct(o)
	o.Structs[key] = s
	o.mu.Unlock()

	return s.Iterate(reflect.ValueOf(ptr).Elem())
}

// Save
// 反射结果保存到JSON文件中.
func (o *Reflection) Save() error {
	for k, v := range o.Structs {
		if err := func(fk string, fp *Struct) error {
			for mk, mp := range map[string]interface{}{
				base.JsonFileItem(fk): fp.Items(),
				base.JsonFileCode(fk): fp.Map(),
			} {
				if err := o.save(mk, mp); err != nil {
					return err
				}
			}
			return nil
		}(k, v); err != nil {
			return err
		}
	}
	return nil
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *Reflection) save(name string, ptr interface{}) error {
	var (
		path     = fmt.Sprintf("%s%s%s/%s", o.BasePath, o.StoragePath, o.TmpPath, name)
		buf, err = json.MarshalIndent(ptr, "", "    ")
	)

	if err != nil {
		return err
	}

	// 1. 检查路径.
	m := regexp.MustCompile(`^(\S+)/([^/]+)$`).FindStringSubmatch(path)
	if len(m) != 3 {
		return fmt.Errorf("invalid file path: %v", m[1])
	}

	// 2. 创建目录.
	if _, ok := o.directories[m[1]]; !ok {
		if err = os.MkdirAll(m[1], os.ModePerm); err != nil {
			return err
		}
		o.directories[m[1]] = true
	}

	// 3. 写入内容.
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	_, err = file.Write(buf)
	return err
}

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
	"github.com/fuyibing/gdoc/config"
	"os"
	"reflect"
	"regexp"
	"strings"
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

// Error
// 打印错误.
func (o *Reflection) Error(err error) {
	println("[REFLECT ERR]", err.Error())
}

// Info
// 打印信息.
func (o *Reflection) Info(format string, args ...interface{}) {
	println("[REFLECTIONS]", fmt.Sprintf(format, args...))
}

// Parse
// 解析结构体.
//
//   ref.Parse("pkg.Example", &pkg.Example{})
func (o *Reflection) Parse(key string, ptr interface{}) {
	o.mu.Lock()

	// 重复检查.
	if _, ok := o.Structs[key]; ok {
		o.mu.Unlock()
		return
	}

	// 解析结果.
	s := NewStruct(o)
	o.Structs[key] = s
	o.mu.Unlock()

	// 迭代过程.
	x := reflect.ValueOf(ptr).Elem()
	s.Iterate(x)
}

// Save
// 反射结果保存到JSON文件中.
func (o *Reflection) Save() {
	for k, v := range o.Structs {
		func(fk string, fp *Struct) {
			for mk, mp := range map[string]interface{}{
				base.JsonFileItem(fk): fp.Items(),
				base.JsonFileCode(fk): fp.Map(),
			} {
				o.save(mk, mp)
			}
		}(k, v)
	}
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *Reflection) save(name string, ptr interface{}) {
	var (
		path     = fmt.Sprintf("%s%s%s/%s", o.BasePath, o.StoragePath, o.TmpPath, name)
		buf, err = json.MarshalIndent(ptr, "", "    ")
		file     *os.File
	)

	// 打印错误.
	defer func() {
		if file != nil {
			err = file.Close()
		}
		if err != nil {
			o.Error(err)
		}
	}()

	// 转化错误.
	if err != nil {
		return
	}

	// 检查路径.
	m := regexp.MustCompile(`^(\S+)/([^/]+)$`).FindStringSubmatch(path)
	if len(m) != 3 {
		err = fmt.Errorf("invalid file path: %v", m[1])
		return
	}

	// 创建目录.
	if _, ok := o.directories[m[1]]; !ok {
		if err = os.MkdirAll(m[1], os.ModePerm); err != nil {
			return
		}
		o.directories[m[1]] = true
	}

	// 打开文件.
	if file, err = os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm); err != nil {
		return
	}

	// 写入内容.
	if _, err = file.Write(buf); err == nil {
		o.Info("create file: %v",
			strings.TrimPrefix(path, config.Path.GetBase()+"/"),
		)
	}
}

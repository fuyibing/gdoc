// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

// Package local
//
// 基于扫描结果, 生成Markdown文件到指定的位置.
package local

import (
	"fmt"
	"github.com/fuyibing/gdoc/config"
	"github.com/fuyibing/gdoc/managers"
	"strings"
)

type (
	// Adapter
	// 适配器接口.
	Adapter interface {
		// Run
		// 执行过程.
		Run() error
	}

	adapter struct {
		manager managers.Management
	}
)

// NewAdapter
// 创建适配器.
func NewAdapter(manager managers.Management) Adapter {
	return &adapter{
		manager: manager,
	}
}

// Run
// 执行过程.
func (o *adapter) Run() error {
	for _, call := range []func() error{
		o.runMethod,
		o.runReadme,
	} {
		if err := call(); err != nil {
			return err
		}
	}
	return nil
}

// 生成文件名.
func (o *adapter) genMethodFilename(requestMethod, requestUrl, prefix string) (s string) {
	s = fmt.Sprintf("%s%s.%s.md", prefix, requestUrl, requestMethod)
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimPrefix(s, ".")
	s = fmt.Sprintf("%s%s/%s", config.Path.GetBase(), config.Path.GetStorage(), strings.ToLower(s))
	return
}

func (o *adapter) runMethod() error {
	for _, c := range o.manager.GetScanner().GetControllers() {
		for _, m := range c.GetMethods() {
			// 忽略.
			if m.GetComment().IsIgnored() {
				continue
			}

			// 导出.
			if err := NewMethodTpl(o, m).Save(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *adapter) runReadme() error { return nil }

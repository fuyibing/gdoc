// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package managers

import (
	"fmt"
	"github.com/fuyibing/gdoc/config"
	"os"
	"regexp"
	"strings"
)

var (
	Manager Management
)

type (
	Management interface {
		GetCompile() Compile
		GetLogger() Logger
		GetScanner() Scanner
		SaveFile(path, text string) error
	}

	management struct {
		Compile Compile
		Logger  Logger
		Scanner Scanner

		makeDirectories map[string]bool
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *management) GetCompile() Compile              { return o.Compile }
func (o *management) GetLogger() Logger                { return o.Logger }
func (o *management) GetScanner() Scanner              { return o.Scanner }
func (o *management) SaveFile(path, text string) error { return o.saveFile(path, text) }

// /////////////////////////////////////////////////////////////
// Access file
// /////////////////////////////////////////////////////////////

func (o *management) saveFile(path, text string) error {
	// 1. 检查路径.
	m := regexp.MustCompile(`^(\S+)/([^/]+)$`).FindStringSubmatch(path)
	if len(m) != 3 {
		return fmt.Errorf("invalid file path: %v", m[1])
	}

	// 2. 创建目录.
	if _, ok := o.makeDirectories[m[1]]; !ok {
		if err := os.MkdirAll(m[1], os.ModePerm); err != nil {
			return err
		}
		o.makeDirectories[m[1]] = true
	}

	// 3. 写入内容.
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	if _, err = file.WriteString(text); err == nil {
		o.Logger.Info("[save] %s", strings.TrimPrefix(path, config.Path.GetBase()))
	}
	return err
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *management) init() *management {
	o.Logger = (&logger{}).init()
	o.Compile = (&compile{Manager: o}).init()
	o.Scanner = (&scanner{Manager: o}).init()

	o.makeDirectories = make(map[string]bool)
	return o
}

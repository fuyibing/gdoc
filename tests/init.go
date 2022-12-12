// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package tests

import (
	"github.com/fuyibing/gdoc/config"
	"sync"
)

var (
	// basePath = "/Users/fuyibing/codes/git.uniondrug.com/gs-fin-es"
	basePath = "/Users/fuyibing/codes/git.uniondrug.com/gs-fin-monitor"
)

func init() {
	new(sync.Once).Do(func() {
		config.Path.SetBase(basePath)
	})
}

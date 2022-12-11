// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package tests

import (
	"github.com/fuyibing/gdoc/config"
	"testing"
)

func TestConfigPath(t *testing.T) {
	t.Logf("base: %s", config.Path.GetBase())
	t.Logf("controller: %s", config.Path.GetController())
	t.Logf("storage: %s", config.Path.GetStorage())
	t.Logf("tmp: %s", config.Path.GetTmp())
}

// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package tests

import (
	"encoding/json"
	"github.com/fuyibing/gdoc/conf"
	"testing"
)

func TestConfig(t *testing.T) {

	buf, _ := json.Marshal(conf.Config)
	t.Logf("config: %s", buf)

}

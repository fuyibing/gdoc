// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package tests

import (
	"github.com/fuyibing/gdoc/conf"
)

func init() {
	// BasePath := "."
	// BasePath := "/Users/fuyibing/go/src/ydyun360-medical-insurance"
	// BasePath := "/Users/fuyibing/codes/git.uniondrug.com/gs-fin-es"
	BasePath := "/Users/fuyibing/codes/git.uniondrug.com/gs-fin-monitor"

	ControllerPath := "/app/controllers"
	// ControllerPath := "/apps/basics/controllers"

	conf.Path.SetBasePath(BasePath)
	conf.Path.SetControllerPath(ControllerPath)
	conf.Config.Load()
}

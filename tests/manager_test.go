// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package tests

import (
	"github.com/fuyibing/gdoc/managers"
	"testing"
)

func TestManagerScanner(t *testing.T) {
	if err := managers.Manager.GetScanner().Scan(); err != nil {
		t.Errorf("scan: %v", err)
		return
	}

	// Compile.
	for _, call := range []func() error{
		managers.Manager.GetCompile().Configure,
		managers.Manager.GetCompile().Make,
		managers.Manager.GetCompile().Install,
	} {
		if err := call(); err != nil {
			t.Errorf("compile: %v", err)
		}
	}

	// for _, c := range managers.Manager.GetScanner().GetControllers() {
	// 	for _, m := range c.GetMethods() {
	// 		if r := m.GetComment().GetRequest(); r != nil {
	// 			t.Logf("request: %d - %v.%v", r.GetLine(), r.GetPackage(), r.GetName())
	// 		}
	// 	}
	// }
}

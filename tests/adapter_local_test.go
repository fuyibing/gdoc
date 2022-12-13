// author: wsfuyibing <websearch@163.com>
// date: 2022-12-13

package tests

import (
	"github.com/fuyibing/gdoc/adapters/local"
	"github.com/fuyibing/gdoc/managers"
	"testing"
)

func TestAdapterLocal(t *testing.T) {
	if err := managers.Manager.GetScanner().Scan(); err != nil {
		t.Errorf("adapter scan: %v", err)
		return
	}

	if err := local.NewAdapter(managers.Manager).Run(); err != nil {
		t.Errorf("adapter result: %v", err)
		return
	}

	t.Logf("adapter complete")
}

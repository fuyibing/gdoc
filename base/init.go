// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Mapper = (&mapping{}).init()
	})
}

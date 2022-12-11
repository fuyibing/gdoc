// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package managers

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Manager = (&management{}).init()
	})
}

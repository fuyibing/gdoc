// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package scanners

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Scanner = (&scanner{}).init()
	})
}

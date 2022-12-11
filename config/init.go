// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package config

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Path = (&path{}).init()
		Regex = (&regex{}).init()
	})
}

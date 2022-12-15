// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package conf

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Config = (&ConfigManager{}).init()
		Debugger = (&DebuggerManager{}).init()
		Path = (&PathManager{}).init()
		Regex = (&RegexManager{}).init()
	})
}

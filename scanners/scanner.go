// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package scanners

import (
	"fmt"
	"github.com/fuyibing/gdoc/conf"
	"os"
)

var (
	Scanner ScannerManager
)

type (
	ScannerManager interface {
		Scan()
	}

	scanner struct {
	}
)

func (o *scanner) Scan() { o.scan("") }

func (o *scanner) init() *scanner {
	return o
}

func (o *scanner) scan(prefix string) {
	var (
		ds   []os.DirEntry
		err  error
		path = fmt.Sprintf("%s%s", conf.Path.GetBasePath(), conf.Path.GetControllerPath())
	)

	if prefix != "" {
		path = fmt.Sprintf("%s%s", path, prefix)
	}

	if ds, err = os.ReadDir(path); err != nil {
		conf.Debugger.Error("[scanner:dir] scan directory: %v", err)
		return
	}

	conf.Debugger.Info("[scanner:dir] scan directory: %s%s", conf.Path.GetControllerPath(), prefix)
	for _, d := range ds {
		// File or directory start with dot should be ignored.
		if conf.Regex.HiddenFile.MatchString(d.Name()) {
			continue
		}

		// Directory recursion.
		if d.IsDir() {
			o.scan(fmt.Sprintf("%s/%s", prefix, d.Name()))
			continue
		}

		// Ignore not source code (.go) files.
		if !conf.Regex.SourceFile.MatchString(d.Name()) {
			continue
		}

		// Found golang source code file.
		(&reader{prefix: prefix, name: d.Name()}).init().run()
	}
}

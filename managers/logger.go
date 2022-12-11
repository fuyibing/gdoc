// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package managers

import (
	"fmt"
	"os"
)

type (
	Logger interface {
		Info(format string, args ...interface{})
		Error(format string, args ...interface{})
	}

	logger struct {
	}
)

// /////////////////////////////////////////////////////////////
// Interface method
// /////////////////////////////////////////////////////////////

func (o *logger) Info(format string, args ...interface{})  { o.println(os.Stdout, format, args...) }
func (o *logger) Error(format string, args ...interface{}) { o.println(os.Stderr, format, args...) }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *logger) println(w *os.File, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, fmt.Sprintf(format, args...)+"\n")
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *logger) init() *logger {
	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package conf

import (
	"fmt"
	"os"
)

var (
	Debugger *DebuggerManager
)

type (
	// DebuggerManager
	// logger progress message.
	DebuggerManager struct{}
)

// Error
// redirect message into standard error file.
func (o *DebuggerManager) Error(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "[ERROR] "+fmt.Sprintf(format, args...)+"\n")
}

// Info
// redirect message into standard output file.
func (o *DebuggerManager) Info(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "[INFO] "+fmt.Sprintf(format, args...)+"\n")
}

// Initialize
// instance fields.
func (o *DebuggerManager) init() *DebuggerManager {
	return o
}

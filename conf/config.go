// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

var (
	Config *ConfigManager
)

type (
	// ConfigManager
	// suit for global configuration.
	ConfigManager struct {
		Deploy      *DeployConfig
		Description string
		Module      string
		Name        string
	}
)

// Load
// configurations from json file.
func (o *ConfigManager) Load() {
	o.load()
	o.loadModule()
	o.defaults()
}

// Fill default
// fields if not specified.
func (o *ConfigManager) defaults() {
	// Init
	// deploy settings.
	if o.Deploy == nil {
		o.Deploy = (&DeployConfig{}).init()
	}

	// Fill
	// deploy default fields.
	o.Deploy.defaults()

	// Fill
	// default module name.
	if o.Module == "" {
		o.Module = "sketch"
	}

	// Fill
	// document name.
	if o.Name == "" {
		o.Name = fmt.Sprintf("sketch, not configured in %s", Path.GetDocumentJsonFile())
	}
}

// Initialize
// config manager.
func (o *ConfigManager) init() *ConfigManager {
	return o
}

// Load fields
// from json file.
func (o *ConfigManager) load() {
	buf, err := os.ReadFile(fmt.Sprintf("%s/%s", Path.GetBasePath(), Path.GetDocumentJsonFile()))

	// Logger
	// load result.
	defer func() {
		if err != nil {
			Debugger.Error("load json config: %v", err)
		} else {
			Debugger.Info("load json config: %s", Path.GetDocumentJsonFile())
		}
	}()

	// Unmarshal
	// into fields.
	if err == nil {
		err = json.Unmarshal(buf, o)
	}
}

// Load module name
// from go.mod file.
func (o *ConfigManager) loadModule() {
	if buf, err := os.ReadFile(fmt.Sprintf("%s/%s", Path.basePath, "go.mod")); err == nil {
		if m := regexp.MustCompile(`\nmodule\s+([^\n]+)\n`).FindStringSubmatch(fmt.Sprintf("\n%s\n", buf)); len(m) == 2 {
			o.Module = m[1]
		}
	}
}

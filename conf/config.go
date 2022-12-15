// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

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
	ConfigManager struct {
		// Document name.
		//
		// Default is sketch.
		Name string

		// Document description.
		Description string

		Module string // sketch

		// Application deploy config
		Deploy *DeployConfig
	}
)

func (o *ConfigManager) Load() {
	o.load()
	o.loadModule()
	o.defaults()
}

func (o *ConfigManager) defaults() {
	if o.Module == "" {
		o.Module = "sketch"
	}

	if o.Name == "" {
		o.Name = fmt.Sprintf("sketch, configured by %s", Path.GetDocumentJsonFile())
	}

	if o.Deploy == nil {
		o.Deploy = (&DeployConfig{}).init()
	}
	o.Deploy.defaults()
}

func (o *ConfigManager) init() *ConfigManager {
	return o
}

func (o *ConfigManager) load() {
	var (
		buf, err = os.ReadFile(fmt.Sprintf("%s/%s", Path.GetBasePath(), Path.GetDocumentJsonFile()))
	)

	defer func() {
		if err != nil {
			Debugger.Error("load json config: %v", err)
		} else {
			Debugger.Info("load json config: %s", Path.GetDocumentJsonFile())
		}
	}()

	if err == nil {
		err = json.Unmarshal(buf, o)
	}
}

func (o *ConfigManager) loadModule() {
	if buf, err := os.ReadFile(fmt.Sprintf("%s/%s", Path.basePath, "go.mod")); err == nil {
		if m := regexp.MustCompile(`\nmodule\s+([^\n]+)\n`).FindStringSubmatch(fmt.Sprintf("\n%s\n", buf)); len(m) == 2 {
			o.Module = m[1]
		}
	}
}

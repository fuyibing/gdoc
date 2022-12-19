// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	Path *PathManager
)

const (
	PathCode = "code"
	PathItem = "item"
)

type (
	// PathManager
	// register environment working path.
	PathManager struct {
		basePath         string
		controllerPath   string
		documentPath     string
		documentJsonFile string
		tmpPath          string

		directories map[string]bool
	}
)

func (o *PathManager) GenerateCodeFile(key string) string { return o.generateFilename(key, PathCode) }
func (o *PathManager) GenerateItemFile(key string) string { return o.generateFilename(key, PathItem) }
func (o *PathManager) GetBasePath() string                { return o.basePath }
func (o *PathManager) GetControllerPath() string          { return o.controllerPath }
func (o *PathManager) GetDocumentJsonFile() string        { return o.documentJsonFile }
func (o *PathManager) GetDocumentPath() string            { return o.documentPath }
func (o *PathManager) GetTmpPath() string                 { return o.tmpPath }
func (o *PathManager) SavePath(path, text string)         { o.saveFile(path, text) }
func (o *PathManager) SetBasePath(s string)               { o.basePath, _ = filepath.Abs(s) }
func (o *PathManager) SetControllerPath(s string)         { o.controllerPath = s }
func (o *PathManager) SetDocumentJsonFile(s string)       { o.documentJsonFile = s }
func (o *PathManager) SetDocumentPath(s string)           { o.documentPath = s }
func (o *PathManager) SetTmpPath(s string)                { o.tmpPath = s }

// Generate
// file name.
//
// - trim prefix
// - change slash as h-line
func (o *PathManager) generateFilename(key, prefix string) (str string) {
	str = strings.TrimPrefix(key, "/")
	str = strings.ReplaceAll(str, "/", "-")
	str = fmt.Sprintf("%s-%s.json", prefix, str)
	return
}

// Initialize
// instance fields.
func (o *PathManager) init() *PathManager {
	o.basePath, _ = filepath.Abs(".")
	o.controllerPath = "/app/controllers"
	o.documentJsonFile = "gdoc.json"
	o.documentPath = "/docs/api"
	o.tmpPath = "/.tmp"

	o.directories = make(map[string]bool)
	return o
}

// Write content
// to specified file.
func (o *PathManager) saveFile(path, text string) {
	var (
		err  error
		file *os.File
	)

	// Called
	// when progress end.
	defer func() {
		// Close file handler
		// if opened.
		if file != nil {
			err = file.Close()
		}

		// Logger
		// file write result.
		if err != nil {
			Debugger.Error("[save] %v", err)
		} else {
			Debugger.Info("[save] %s", strings.TrimPrefix(path, Path.GetBasePath()))
		}
	}()

	// Make directory
	// when called at first time.
	if m := regexp.MustCompile(`^(.+)/([^/]+)$`).FindStringSubmatch(path); len(m) == 3 {
		if _, ok := o.directories[m[1]]; !ok {
			if err = os.MkdirAll(m[1], os.ModePerm); err != nil {
				return
			}
		}
	}

	// Open file.
	//
	// - truncate content if exists
	// - create file if not exists
	if file, err = os.OpenFile(path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm); err != nil {
		return
	}

	// Write progress.
	_, err = file.WriteString(text)
}

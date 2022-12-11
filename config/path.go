// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	Path PathManager

	pathDefaultBase       = "./"
	pathDefaultConfigJson = "gdoc.json"
	pathDefaultController = "/app/controllers"
	pathDefaultStorage    = "/docs/api"
	pathDefaultTmp        = "/.tmp"

	pathRegexReplaceSlashes = regexp.MustCompile(`/+`)
)

type (
	PathManager interface {
		GetBase() string
		GetConfigJson() string
		GetController() string
		GetStorage() string
		GetTmp() string
		SetBase(s string)
		SetConfigJson(s string)
		SetController(s string)
		SetStorage(s string)
		SetTmp(s string)
	}

	path struct {
		BasePath       string
		ConfigJson     string
		ControllerPath string
		StoragePath    string
		TmpPath        string
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *path) GetBase() string        { return o.BasePath }
func (o *path) GetConfigJson() string  { return o.ConfigJson }
func (o *path) GetController() string  { return o.ControllerPath }
func (o *path) GetStorage() string     { return o.StoragePath }
func (o *path) GetTmp() string         { return o.TmpPath }
func (o *path) SetBase(s string)       { o.BasePath = o.abs(s) }
func (o *path) SetConfigJson(s string) { o.ConfigJson = s }
func (o *path) SetController(s string) { o.ControllerPath = o.rel(s) }
func (o *path) SetStorage(s string)    { o.StoragePath = o.rel(s) }
func (o *path) SetTmp(s string)        { o.TmpPath = o.rel(s) }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *path) abs(s string) string {
	src, _ := filepath.Abs(s)
	return src
}

func (o *path) rel(s string) string {
	s = pathRegexReplaceSlashes.ReplaceAllString(s, "/")
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimSuffix(s, "/")

	if s != "" {
		s = fmt.Sprintf("/%s", s)
	}

	return s
}

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *path) init() *path {
	o.BasePath = o.abs(pathDefaultBase)
	o.ConfigJson = pathDefaultConfigJson
	o.ControllerPath = pathDefaultController
	o.StoragePath = pathDefaultStorage
	o.TmpPath = pathDefaultTmp
	return o
}

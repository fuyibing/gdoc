// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package conf

import (
	"regexp"
)

var (
	// Regex
	// 正则单例.
	Regex *RegexManager
)

type (
	RegexManager struct {
		HiddenFile *regexp.Regexp

		SourceAnnotation      *regexp.Regexp
		SourceComment         *regexp.Regexp
		SourceController      *regexp.Regexp
		SourceControllerGroup *regexp.Regexp
		SourceFile            *regexp.Regexp
		SourceMethod          *regexp.Regexp
	}
)

func (o *RegexManager) init() *RegexManager {
	o.HiddenFile = regexHiddenFile

	o.SourceAnnotation = regexSourceAnnotation
	o.SourceComment = regexSourceComment
	o.SourceController = regexSourceController
	o.SourceControllerGroup = regexSourceControllerGroup
	o.SourceFile = regexSourceFile
	o.SourceMethod = regexSourceMethod
	return o
}

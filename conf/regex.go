// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package conf

import (
	"regexp"
)

var (
	Regex *RegexManager
)

type (
	// RegexManager
	// register regular expression.
	RegexManager struct {
		HiddenFile *regexp.Regexp

		SourceAnnotation, SourceAnnotationSimple *regexp.Regexp
		SourceComment                            *regexp.Regexp
		SourceController, SourceControllerGroup  *regexp.Regexp
		SourceFile                               *regexp.Regexp
		SourceMethod                             *regexp.Regexp
		SourceRoute, SourceRouteMethod           *regexp.Regexp
	}
)

// Initialize instance field
// with default regular expressions.
func (o *RegexManager) init() *RegexManager {
	o.HiddenFile = regexHiddenFile

	o.SourceAnnotation = regexSourceAnnotation
	o.SourceAnnotationSimple = regexSourceAnnotationSimple
	o.SourceComment = regexSourceComment
	o.SourceController = regexSourceController
	o.SourceControllerGroup = regexSourceControllerGroup
	o.SourceFile = regexSourceFile
	o.SourceMethod = regexSourceMethod
	o.SourceRoute = regexSourceRoute
	o.SourceRouteMethod = regexSourceRouteMethod
	return o
}

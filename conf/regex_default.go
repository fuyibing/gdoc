// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package conf

import (
	"regexp"
)

var (
	regexHiddenFile = regexp.MustCompile(`^\.`)

	regexSourceAnnotation      = regexp.MustCompile(`^\s*@([A-Z][_a-zA-Z0-9]*)\s*\(([^)]*)\)`)
	regexSourceFile            = regexp.MustCompile(`\.go$`)
	regexSourceComment         = regexp.MustCompile(`^/{2,}\s?([^\n]*)`)
	regexSourceController      = regexp.MustCompile(`type\s+([_a-zA-Z0-9]*Controller)\s*struct\s*\{`)
	regexSourceControllerGroup = regexp.MustCompile(`^([_a-zA-Z0-9]*Controller)\s*struct\s*\{`)
	regexSourceMethod          = regexp.MustCompile(`^func\s*\(\s*[_a-zA-Z0-9]*\s*[*]?([_a-zA-Z0-9]*Controller)\s*\)\s*([A-Z][_a-zA-Z0-9]*)\s*\(.*\)`)
)

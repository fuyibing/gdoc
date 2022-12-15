// author: wsfuyibing <websearch@163.com>
// date: 2022-12-15

package tpl

import (
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"regexp"
	"strings"
)

func linkName(m base.Method) (str string) {
	str = fmt.Sprintf("%s.%s.md", m.GetRequestUrl(), m.GetRequestMethod())
	str = strings.ToLower(str)
	str = regexp.MustCompile(`^/\.`).ReplaceAllString(str, "/")
	str = fmt.Sprintf("/api%s", str)
	return
}

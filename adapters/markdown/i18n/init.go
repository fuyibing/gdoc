// author: wsfuyibing <websearch@163.com>
// date: 2022-12-15

package i18n

import (
	"sync"
)

type (
	Language map[string]string
)

var (
	defaultLanguage Language
)

func Lang(key string) string {
	if s, ok := defaultLanguage[key]; ok {
		return s
	}
	return key
}

func SetLang(lang Language) {
	defaultLanguage = lang
}

func init() {
	new(sync.Once).Do(func() {
		defaultLanguage = zh
	})
}

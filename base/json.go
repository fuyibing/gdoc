// author: wsfuyibing <websearch@163.com>
// date: 2022-12-12

package base

import (
	"fmt"
	"strings"
)

func JsonFilename(key string) string {
	return strings.ToLower(strings.ReplaceAll(key, "/", "_"))
}

func JsonFileCode(key string) string {
	return fmt.Sprintf("%s.code.json", JsonFilename(key))
}

func JsonFileItem(key string) string {
	return fmt.Sprintf("%s.item.json", JsonFilename(key))
}

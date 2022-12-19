// author: wsfuyibing <websearch@163.com>
// date: 2020-12-13

package tpl

import (
	"crypto/md5"
	"fmt"
	"github.com/fuyibing/gdoc/conf"
)

type (
	// Info
	// 基础资料.
	Info struct {
		Id          string `json:"_postman_id"`
		Description string `json:"description"`
		Name        string `json:"name"`
		Schema      string `json:"schema"`
	}
)

func NewInfo() *Info {
	return (&Info{
		Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
	}).init()
}

func (o *Info) init() *Info {
	s := fmt.Sprintf("%x", md5.Sum([]byte(conf.Config.Module)))
	o.Id = fmt.Sprintf("%s-%s-%s-%s", s[0:8], s[8:12], s[12:16], s[16:])
	return o
}

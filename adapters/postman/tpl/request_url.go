// author: wsfuyibing <websearch@163.com>
// date: 2022-12-15

package tpl

import (
	"fmt"
	"strings"
)

type (
	RequestUrl struct {
		Host     []string `json:"host"`
		Path     []string `json:"path"`
		Query    []string `json:"query,omitempty"`
		Protocol string   `json:"protocol"`
		Raw      string   `json:"raw"`
	}
)

func NewRequestUrl() *RequestUrl {
	return (&RequestUrl{}).init()
}

func (o *RequestUrl) End() {
	o.Raw = fmt.Sprintf("%s://%s/%s", o.Protocol, strings.Join(o.Host, "."), strings.Join(o.Path, "/"))
}

func (o *RequestUrl) SetHost(protocol, host string, port int) {
	o.Protocol = protocol

	if port > 0 {
		s := strings.ToLower(protocol)

		if (s == "http" && port == 80) || (s == "https" && port == 443) {
			o.Host = []string{host}
			return
		}
	}

	o.Host = []string{fmt.Sprintf("%s:%d", host, port)}
}

func (o *RequestUrl) SetPath(path string) {
	strings.TrimPrefix(path, "/")
	o.Path = strings.Split(path, "/")
}

func (o *RequestUrl) init() *RequestUrl {
	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package conf

import (
	"fmt"
)

type (
	DeployConfig struct {
		Protocol string
		Host     string
		Port     int
	}
)

func (o *DeployConfig) Full() string {
	if (o.Protocol == "http" && o.Port == 80) || (o.Protocol == "https" && o.Port == 443) {
		return fmt.Sprintf("%s://%s", o.Protocol, o.Host)
	}
	return fmt.Sprintf("%s://%s:%d", o.Protocol, o.Host, o.Port)
}

func (o *DeployConfig) defaults() {
	if o.Protocol == "" {
		o.Protocol = "http"
	}
	if o.Host == "" {
		o.Host = "127.0.0.1"
	}
	if o.Port == 0 {
		o.Port = 8080
	}
}

func (o *DeployConfig) init() *DeployConfig {
	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package conf

type (
	DeployConfig struct {
		Protocol string
		Host     string
		Port     int
	}
)

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

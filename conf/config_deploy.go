// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package conf

import (
	"fmt"
)

type (
	// DeployConfig
	// for application deploy settings.
	DeployConfig struct {
		Protocol string
		Host     string
		Port     int
	}
)

// Full
// return host and port with protocol.
func (o *DeployConfig) Full() string {
	// Ignore port
	// if matched with default protocol.
	if (o.Protocol == "http" && o.Port == 80) || (o.Protocol == "https" && o.Port == 443) {
		return fmt.Sprintf("%s://%s", o.Protocol, o.Host)
	}

	// Return
	// standard string.
	return fmt.Sprintf("%s://%s:%d", o.Protocol, o.Host, o.Port)
}

// Fill
// default fields.
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

// Initialize
// instance fields.
func (o *DeployConfig) init() *DeployConfig {
	return o
}

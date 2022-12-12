// author: wsfuyibing <websearch@163.com>
// date: 2022-12-11

package base

type Item struct {
	Array       bool        `json:"array"`
	Condition   string      `json:"condition"`
	Description string      `json:"description"`
	Ignored     bool        `json:"ignored"`
	Key         string      `json:"key"`
	Kind        int         `json:"kind"`
	Label       string      `json:"label"`
	Name        string      `json:"name"`
	Required    bool        `json:"required"`
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
	Children    []*Item     `json:"children,omitempty"`
}

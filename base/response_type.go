// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package base

import (
	"encoding/json"
)

type (
	ResponseType int
)

const (
	_ ResponseType = iota

	ResponseData
	ResponseList
	ResponsePaging
)

func (a ResponseType) Render(b []byte) string {
	var (
		buf []byte
		err error
		res = make(map[string]interface{})
		v   = NewResponseTpl()
	)

	// Return
	// empty string.
	if b == nil {
		return ""
	}

	// Return
	// error string.
	if err = json.Unmarshal(b, &res); err != nil {
		return err.Error()
	}

	switch a {
	case ResponseList:
		v.AsList(res)

	case ResponsePaging:
		v.AsPaging(res)

	default:
		v.AsData(res)
	}

	// Marshal
	// into json string.
	if buf, err = json.MarshalIndent(v, "", "    "); err != nil {
		return err.Error()
	}

	// Return
	// json string.
	return string(buf)
}

func (a ResponseType) toData(body []byte) {
}

func (a ResponseType) toList(body []byte) {}

func (a ResponseType) toPaging(body []byte) {}

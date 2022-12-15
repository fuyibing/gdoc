// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

type (
	ResponseType int
)

const (
	_ ResponseType = iota

	ResponseData
	ResponseList
	ResponsePaging
)

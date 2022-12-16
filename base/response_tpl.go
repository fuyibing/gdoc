// author: wsfuyibing <websearch@163.com>
// date: 2022-12-16

package base

type (
	ResponseTpl struct {
		Data     interface{} `json:"data"`
		DataType string      `json:"dataType"`
		Errno    int         `json:"errno"`
		Error    string      `json:"error"`
	}

	ResponseTplPaging struct {
		First      int   `json:"first"`
		Before     int   `json:"before"`
		Current    int   `json:"current"`
		Next       int   `json:"next"`
		Last       int   `json:"last"`
		Limit      int   `json:"limit"`
		TotalPages int   `json:"totalPages"`
		TotalItems int64 `json:"totalItems"`
	}
)

func NewResponseTpl() *ResponseTpl {
	return &ResponseTpl{}
}

func (o *ResponseTpl) AsData(v interface{}) {
	o.Data = v
	o.DataType = "OBJECT"
}

func (o *ResponseTpl) AsList(v interface{}) {
	o.Data = []interface{}{v}
	o.DataType = "LIST"
}

func (o *ResponseTpl) AsPaging(v interface{}) {
	o.Data = map[string]interface{}{
		"body":   []interface{}{v},
		"paging": (&ResponseTplPaging{}).init(),
	}
	o.DataType = "PAGING"
}

func (o *ResponseTplPaging) init() *ResponseTplPaging {
	o.First = 1
	o.Before = 1
	o.Current = 1
	o.Next = 1
	o.Last = 1
	o.Limit = 10
	o.TotalPages = 1
	o.TotalItems = 1
	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

type (
	Method interface {
		GetComment() *Comment
		GetController() Controller
	}

	method struct {
		comment    *Comment
		controller Controller
	}
)

func NewMethod(controller Controller) Method {
	return (&method{
		comment:    NewComment(),
		controller: controller,
	}).init()
}

func (o *method) GetComment() *Comment      { return o.comment }
func (o *method) GetController() Controller { return o.controller }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *method) init() *method {
	return o
}

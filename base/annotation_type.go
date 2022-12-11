// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package base

type (
	AnnotationType string
)

const (
	AnnotationError          AnnotationType = "Error"
	AnnotationRequest        AnnotationType = "Request"
	AnnotationResponse       AnnotationType = "Response"
	AnnotationResponseData   AnnotationType = "ResponseData"
	AnnotationResponseList   AnnotationType = "ResponseList"
	AnnotationResponsePaging AnnotationType = "ResponsePaging"

	// AnnotationRoutePrefix
	// use to define on controller.
	//
	//   @RoutePrefix(/)
	//   @RoutePrefix(/example)
	AnnotationRoutePrefix AnnotationType = "RoutePrefix"
)

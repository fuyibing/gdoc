// author: wsfuyibing <websearch@163.com>
// date: 2022-12-14

package base

type (
	AnnotationType string
)

const (
	AnnotationError          AnnotationType = "Error"
	AnnotationHeader         AnnotationType = "Header"
	AnnotationIgnore         AnnotationType = "Ignore"
	AnnotationRequest        AnnotationType = "Request"
	AnnotationResponse       AnnotationType = "Response"
	AnnotationResponseData   AnnotationType = "ResponseData"
	AnnotationResponseList   AnnotationType = "ResponseList"
	AnnotationResponsePaging AnnotationType = "ResponsePaging"
	AnnotationRoutePrefix    AnnotationType = "RoutePrefix"
	AnnotationVersion        AnnotationType = "Version"
)

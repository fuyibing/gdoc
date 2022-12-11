// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package config

import (
	"regexp"
)

var (
	Regex RegexManager

	regexCommentEnd  = regexp.MustCompile(`\.+$`)
	regexCommentLine = regexp.MustCompile(`^/+\s?(.*)$`)

	regexController        = regexp.MustCompile(`(\S*Controller)\s+struct\s*\{`)
	regexMethod            = regexp.MustCompile(`^func\s+\([a-zA-Z0-9]*\s*[*]*(\S*Controller)\s*\)\s+([A-Z][_a-zA-Z0-9]*)\s*\([^)]*\)\s+`)
	regexStructWithPackage = regexp.MustCompile(`^([^.]+)\.([A-Z][_a-zA-Z0-9]*)$`)

	regexRouteMethod = regexp.MustCompile(`^(Get|Post|Head|Options|Patch|Put|Delete)$`)
	regexRouteUrl    = regexp.MustCompile(`^(Get|Post|Head|Options|Patch|Put|Delete)([A-Z][_a-zA-Z0-9]*)`)

	regexAnnotationSimple = regexp.MustCompile(`^\s*@([_a-zA-Z0-9]+)$`)
	regexAnnotationParams = regexp.MustCompile(`^\s*@([_a-zA-Z0-9]+)\s*\(([^)]*)\)$`)

	regexExported   = regexp.MustCompile(`^[A-Z]`)
	regexHiddenFile = regexp.MustCompile(`^\.`)
	regexSourceFile = regexp.MustCompile(`\.go$`)
)

type (
	RegexManager interface {
		GetAnnotationParams() *regexp.Regexp
		GetAnnotationSimple() *regexp.Regexp
		GetCommentEnd() *regexp.Regexp
		GetCommentLine() *regexp.Regexp
		GetController() *regexp.Regexp
		GetExported() *regexp.Regexp
		GetHiddenFile() *regexp.Regexp
		GetMethod() *regexp.Regexp
		GetRouteMethod() *regexp.Regexp
		GetRouteUrl() *regexp.Regexp
		GetSourceFile() *regexp.Regexp
		GetStructWithPackage() *regexp.Regexp
	}

	regex struct {
		commentEnd  *regexp.Regexp
		commentLine *regexp.Regexp

		controller        *regexp.Regexp
		method            *regexp.Regexp
		structWithPackage *regexp.Regexp

		routeMethod *regexp.Regexp
		routeUrl    *regexp.Regexp

		annotationParams *regexp.Regexp
		annotationSimple *regexp.Regexp

		exported   *regexp.Regexp
		hiddenFile *regexp.Regexp
		sourceFile *regexp.Regexp
	}
)

// /////////////////////////////////////////////////////////////
// Interface methods
// /////////////////////////////////////////////////////////////

func (o *regex) GetAnnotationParams() *regexp.Regexp  { return o.annotationParams }
func (o *regex) GetAnnotationSimple() *regexp.Regexp  { return o.annotationSimple }
func (o *regex) GetCommentEnd() *regexp.Regexp        { return o.commentEnd }
func (o *regex) GetCommentLine() *regexp.Regexp       { return o.commentLine }
func (o *regex) GetController() *regexp.Regexp        { return o.controller }
func (o *regex) GetExported() *regexp.Regexp          { return o.exported }
func (o *regex) GetHiddenFile() *regexp.Regexp        { return o.hiddenFile }
func (o *regex) GetMethod() *regexp.Regexp            { return o.method }
func (o *regex) GetRouteMethod() *regexp.Regexp       { return o.routeMethod }
func (o *regex) GetRouteUrl() *regexp.Regexp          { return o.routeUrl }
func (o *regex) GetSourceFile() *regexp.Regexp        { return o.sourceFile }
func (o *regex) GetStructWithPackage() *regexp.Regexp { return o.structWithPackage }

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

// /////////////////////////////////////////////////////////////
// Initialize
// /////////////////////////////////////////////////////////////

func (o *regex) init() *regex {
	o.commentEnd = regexCommentEnd
	o.commentLine = regexCommentLine

	o.controller = regexController
	o.method = regexMethod
	o.structWithPackage = regexStructWithPackage

	o.routeMethod = regexRouteMethod
	o.routeUrl = regexRouteUrl

	o.annotationSimple = regexAnnotationSimple
	o.annotationParams = regexAnnotationParams

	o.exported = regexExported
	o.hiddenFile = regexHiddenFile
	o.sourceFile = regexSourceFile
	return o
}

// author: wsfuyibing <websearch@163.com>
// date: 2020-12-14

package scanners

import (
	"bufio"
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/conf"
	"io"
	"os"
	"regexp"
	"strings"
)

type (
	reader struct {
		path         string
		prefix, name string

		comments map[int]string
	}
)

func (o *reader) findComment(line int, text string) {
	o.comments[line] = text
}

func (o *reader) findController(line int, code, cn string) {
	c := o.genController(cn)
	c.SetPrefix(o.prefix)

	o.genComment(line, code, cn, c.GetComment(), c)
}

func (o *reader) findMethod(line int, code, cn, mn string) {
	var (
		c = o.genController(cn)
		m = c.GetMethod(mn)
	)

	// Create method if not exists.
	if m == nil {
		conf.Debugger.Info("[scanner:method] <%s.%s>", cn, mn)
		m = base.NewMethod(c, mn)
		c.SetMethod(mn, m)
	}

	// Generate comment.
	o.genComment(line, code, mn, m.GetComment(), c)
}

func (o *reader) genAnnotation(comment *base.Comment, controller base.Controller, annotation base.Annotation) {
	switch annotation.GetType() {
	case base.AnnotationRoutePrefix:
		controller.SetPrefix(annotation.GetFirst())

	default:
		comment.AddAnnotation(annotation)
	}
}

func (o *reader) genComment(line int, code, name string, comment *base.Comment, controller base.Controller) {
	// Source code info.
	comment.Path = fmt.Sprintf("%s%s", conf.Path.GetControllerPath(), o.path)
	comment.Line = line
	comment.Code = code

	// Package and name.
	comment.Name = name
	comment.Pkg = fmt.Sprintf("%s%s%s", conf.Config.Module, conf.Path.GetControllerPath(), o.prefix)

	// Comment from collected.
	cs := make([]string, 0)
	for i := line - 1; i > 0; i-- {
		if s, ok := o.comments[i]; ok {
			cs = append(cs, s)
			continue
		}
		break
	}

	// Comment send list.
	us := make([]string, 0)
	for i := len(cs) - 1; i >= 0; i-- {
		us = append(us, cs[i])
	}

	// Comment Parse
	o.genCommentParse(comment, controller, us)
}

func (o *reader) genCommentParse(comment *base.Comment, controller base.Controller, cs []string) {
	for _, s := range cs {
		// Annotation simple.
		if m := conf.Regex.SourceAnnotationSimple.FindStringSubmatch(s); len(m) == 2 {
			o.genAnnotation(comment, controller, base.NewAnnotation(m[1], ""))
			continue
		}

		// Annotation standard.
		if m := conf.Regex.SourceAnnotation.FindStringSubmatch(s); len(m) == 3 {
			o.genAnnotation(comment, controller, base.NewAnnotation(m[1], m[2]))
			continue
		}

		// Send to comment.
		comment.AddText(s)
	}
}

func (o *reader) genController(cn string) base.Controller {
	var (
		k = o.genControllerKey(cn)
		c = base.Mapper.GetController(k)
	)

	// Create controller if not found.
	if c == nil {
		conf.Debugger.Info("[scanner:controller] <%s>", cn)

		c = base.NewController(base.Mapper)
		base.Mapper.SetController(k, c)
	}

	return c
}

func (o *reader) genControllerKey(cn string) string {
	return fmt.Sprintf("%s:%s", o.path, cn)
}

func (o *reader) init() *reader {
	o.comments = make(map[int]string)

	o.path = fmt.Sprintf("%s/%s", o.prefix, o.name)
	return o
}

func (o *reader) run() {
	var (
		err  error
		file *os.File
		src  = fmt.Sprintf("%s%s%s", conf.Path.GetBasePath(), conf.Path.GetControllerPath(), o.path)
	)

	defer func() {
		if file != nil {
			err = file.Close()
		}
		if err != nil {
			conf.Debugger.Error("[scanner:file] read source file: %v", err)
		}
	}()

	if file, err = os.OpenFile(src, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}

	var (
		delim = byte('\n')
		line  = 0
		read  = bufio.NewReader(file)
		text  string
	)

	conf.Debugger.Info("[scanner:file] read source file: %s%s/%s", conf.Path.GetControllerPath(), o.prefix, o.name)

	for {
		line++

		// Read text line by line.
		if text, err = read.ReadString(delim); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		// Ignore empty line.
		if text = strings.TrimSpace(text); text == "" {
			continue
		}

		// Comment line.
		if m := conf.Regex.SourceComment.FindStringSubmatch(text); len(m) == 2 {
			o.findComment(line, m[1])
			continue
		}

		// Controller.
		for _, r := range []*regexp.Regexp{
			conf.Regex.SourceController,
			conf.Regex.SourceControllerGroup,
		} {
			if m := r.FindStringSubmatch(text); len(m) == 2 {
				o.findController(line, text, m[1])
				continue
			}
		}

		// Method.
		if m := conf.Regex.SourceMethod.FindStringSubmatch(text); len(m) == 3 {
			o.findMethod(line, text, m[1], m[2])
			continue
		}
	}
}

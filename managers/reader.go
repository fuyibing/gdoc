// author: wsfuyibing <websearch@163.com>
// date: 2022-12-09

package managers

import (
	"bufio"
	"fmt"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/config"
	"io"
	"os"
	"strings"
)

type (
	// Reader
	// use to read source file into controller.
	Reader interface {
		Run() error
	}

	reader struct {
		Manager Management

		commentMapper map[int]string
		folder, name  string
	}

	readerComment struct {
		Line int
		Text string
	}
)

func NewReader(manager Management, folder, name string) Reader {
	return (&reader{
		Manager: manager,
		folder:  folder, name: name,
	}).init()
}

func (o *reader) Run() (err error) {
	var (
		file *os.File
		path = fmt.Sprintf("%s%s%s/%s",
			config.Path.GetBase(),
			config.Path.GetController(),
			o.folder,
			o.name,
		)
	)

	// Close file resource
	// when end if opened.
	defer func() {
		if file != nil {
			err = file.Close()
		}
	}()

	// Return error
	// if file opened fail.
	if file, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}

	// Prepare
	// read line by line.
	var (
		buf   = bufio.NewReader(file)
		delim = byte('\n')
		line  int
		text  string
	)

	for {
		line++

		// Return or break
		// if error occurred.
		if text, err = buf.ReadString(delim); err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		// Check line text.
		if text = strings.TrimSpace(text); text != "" {
			o.check(line, text)
		}
	}
	return
}

func (o *reader) check(line int, text string) {
	// Comment.
	if m := config.Regex.GetCommentLine().FindStringSubmatch(text); len(m) == 2 {
		o.commentMapper[line] = m[1]
		return
	}

	// Controller.
	if m := config.Regex.GetController().FindStringSubmatch(text); len(m) == 2 {
		if config.Regex.GetExported().MatchString(m[1]) {
			o.checkController(line, m[1])
		}
		return
	}

	// Method.
	if m := config.Regex.GetMethod().FindStringSubmatch(text); len(m) == 3 {
		if config.Regex.GetExported().MatchString(m[1]) && config.Regex.GetExported().MatchString(m[2]) {
			o.checkMethod(line, m[1], m[2])
		}
	}
}

func (o *reader) checkController(line int, cn string) {
	controller := o.doController(cn)
	o.doComment(controller, controller.GetComment(), line)
}

func (o *reader) checkMethod(line int, cn, mn string) {
	var (
		controller = o.doController(cn)
		method     = controller.GetMethod(mn)
	)

	if method == nil {
		method = base.NewMethod(controller, mn)
		controller.SetMethod(mn, method)
	}
	o.doComment(method.GetController(), method.GetComment(), line)
}

func (o *reader) doAnnotation(controller base.Controller, comment base.Comment, line int, key, value string) {
	an := base.NewAnnotation(line, key, value)

	switch an.GetType() {
	case base.AnnotationRequest:
		if s := an.GetFirst(); s != "" {
			comment.SetRequest(line, an.GetFirst())
		}

	case base.AnnotationResponse, base.AnnotationResponseData:
		if s := an.GetFirst(); s != "" {
			comment.SetResponse(base.ResponseData, line, an.GetFirst())
		}

	case base.AnnotationResponseList:
		if s := an.GetFirst(); s != "" {
			comment.SetResponse(base.ResponseList, line, an.GetFirst())
		}

	case base.AnnotationResponsePaging:
		if s := an.GetFirst(); s != "" {
			comment.SetResponse(base.ResponsePaging, line, an.GetFirst())
		}

	case base.AnnotationRoutePrefix:
		if controller != nil {
			controller.SetPrefix(an.GetFirst())
		}
	}
}

func (o *reader) doComment(controller base.Controller, comment base.Comment, line int) {
	cs := make([]readerComment, 0)

	// Group comments
	for i := line - 1; i > 0; i-- {
		if text, ok := o.commentMapper[i]; ok {
			cs = append(cs, readerComment{Line: i, Text: text})
			continue
		}
		break
	}

	// Return if no comment.
	if cn := len(cs); cn > 0 {
		for i := cn - 1; i >= 0; i-- {
			c := cs[i]

			// Annotation on simple.
			if m := config.Regex.GetAnnotationSimple().FindStringSubmatch(c.Text); len(m) == 2 {
				o.doAnnotation(controller, comment, c.Line, m[1], "")
				continue
			}

			// Annotation on standard with params.
			if m := config.Regex.GetAnnotationParams().FindStringSubmatch(c.Text); len(m) == 3 {
				o.doAnnotation(controller, comment, c.Line, m[1], m[2])
				continue
			}

			// Text comment.
			comment.AddText(c.Line, c.Text)
		}
	}
}

func (o *reader) doController(cn string) base.Controller {
	var (
		controller base.Controller
		key        = o.key(cn)
	)

	// Read controller
	// and create if not found.
	if controller = o.Manager.GetScanner().GetController(key); controller == nil {
		controller = base.NewController(cn, o.folder)
		o.Manager.GetScanner().SetController(key, controller)
	}

	return controller
}

func (o *reader) init() *reader {
	o.commentMapper = make(map[int]string)
	return o
}

func (o *reader) key(cn string) string {
	return fmt.Sprintf("%s.%s", o.folder, cn)
}

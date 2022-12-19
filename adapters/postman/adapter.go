// author: wsfuyibing <websearch@163.com>
// date: 2020-12-15

package postman

import (
	"encoding/json"
	"fmt"
	"github.com/fuyibing/gdoc/adapters/postman/tpl"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/conf"
	"sort"
)

type (
	// Adapter
	// use to create postman.json file.
	Adapter struct {
		mapping base.Mapping

		p *tpl.Postman
	}
)

// New
// return adapter instance.
func New(mapping base.Mapping) *Adapter {
	return (&Adapter{
		mapping: mapping,
		p:       tpl.New(),
	}).init()
}

// Run
// adapter instance.
func (o *Adapter) Run() {
	o.runInfo()
	o.runItem()
	o.save()
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

func (o *Adapter) doController(c base.Controller) {
	item := tpl.NewItem()
	list := make([]string, 0)
	maps := make(map[string]base.Method)

	// Sort methods by key.
	for _, m := range c.GetMethods() {
		k := m.GetSortKey()
		list = append(list, k)
		maps[k] = m
	}

	// Iterate by sorted.
	sort.Strings(list)
	for _, k := range list {
		m := maps[k]
		if m.GetComment().Ignored {
			continue
		}
		o.doMethod(item, m)
	}

	// Add to parent.
	if len(item.Item) > 0 {
		item.Name = c.GetComment().GetTitle()
		item.Description = c.GetComment().GetDescription()
		o.p.Add(item)
	}
}

func (o *Adapter) doMethod(parent *tpl.Item, m base.Method) {
	if m.GetComment().Ignored {
		return
	}

	// Basic.
	item := tpl.NewItem()
	item.Name = m.GetComment().GetTitle()
	defer parent.Add(item)

	// Request.
	o.doMethodRequest(item, m)

	// Response.
	for i, r := range m.GetComment().Responses {
		o.doMethodResponse(item, m, i, r)
	}
}

func (o *Adapter) doMethodRequest(item *tpl.Item, m base.Method) {
	request := tpl.NewRequest()
	request.Description = m.GetComment().GetDescription()
	request.Method = m.GetRequestMethod()
	item.Request = request

	o.doMethodRequestBody(request, m)
	o.doMethodRequestHeader(request, m)
	o.doMethodRequestUrl(request, m)
}

func (o *Adapter) doMethodRequestBody(request *tpl.Request, m base.Method) {
	if r := m.GetComment().Request; r != nil {
		conf.Debugger.Info("[json] %s", r.Key)

		x := tpl.NewRequestBody()
		x.Mode = tpl.RequestModeRaw
		x.End(r.Key)

		request.Body = x
	}
}

func (o *Adapter) doMethodRequestHeader(request *tpl.Request, m base.Method) {
	for _, x := range m.GetComment().Headers {
		h := tpl.NewHeader()
		h.Key = x.Key
		h.Value = x.Value
		request.Header = append(request.Header, h)
	}
}

func (o *Adapter) doMethodRequestUrl(request *tpl.Request, m base.Method) {
	url := tpl.NewRequestUrl()
	url.SetHost(conf.Config.Deploy.Protocol, conf.Config.Deploy.Host, conf.Config.Deploy.Port)
	url.SetPath(m.GetRequestUrl())
	url.End()

	request.Url = url
}

func (o *Adapter) doMethodResponse(item *tpl.Item, m base.Method, i int, r *base.Response) {
	response := tpl.NewResponse()
	response.Name = fmt.Sprintf("Response #%d", i+1)
	response.Request = item.Request
	response.End(r.Key)
	item.SetResponse(response)
}

func (o *Adapter) init() *Adapter {
	return o
}

func (o *Adapter) runInfo() {
	x := tpl.NewInfo()
	x.Name = conf.Config.Name
	x.Description = conf.Config.Description
	o.p.SetInfo(x)
}

func (o *Adapter) runItem() {
	var (
		list = make([]string, 0)
		maps = make(map[string]base.Controller)
	)

	for _, c := range o.mapping.GetControllers() {
		list = append(list, c.GetSortKey())
		maps[c.GetSortKey()] = c
	}

	sort.Strings(list)

	for _, k := range list {
		o.doController(maps[k])
	}
}

func (o *Adapter) save() {
	buf, err := json.MarshalIndent(o.p, "", "    ")
	if err != nil {
		conf.Debugger.Error("[adapter=postman] %v", err)
		return
	}

	conf.Path.SavePath(fmt.Sprintf(
		"%s%s/postman.json",
		conf.Path.GetBasePath(),
		conf.Path.GetDocumentPath(),
	), string(buf))
}

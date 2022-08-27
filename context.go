package THz

import (
	"math"

	"github.com/valyala/fasthttp"
	"go.x2ox.com/THz/render"
	"go.x2ox.com/sorbifolia/httprouter"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type Handler httprouter.Handler[Context]

type Context struct {
	fc  *fasthttp.RequestCtx
	thz *THz

	params   httprouter.Params
	handlers httprouter.Handlers[Context]
	index    int
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Request() *fasthttp.Request   { return &c.fc.Request }
func (c *Context) Response() *fasthttp.Response { return &c.fc.Response }

func (c *Context) SetHeader(k, v string) *Context   { c.fc.Response.Header.Set(k, v); return c }
func (c *Context) SetLocation(v string) *Context    { return c.SetHeader("Location", v) }
func (c *Context) SetContentType(v string) *Context { return c.SetHeader("Content-Type", v) }

func (c *Context) Status(code int) *Context { c.fc.Response.SetStatusCode(code); return c }
func (c *Context) Abort() *Context          { c.index = math.MaxInt - 1; return c }
func (c *Context) Render(r render.Render) {
	r.WriteContentType(c.fc)
	if err := r.Render(c.fc); err != nil {
		panic(err)
	}
}

func (c *Context) JSON(data any)               { c.Render(render.JSON{Data: data}) }
func (c *Context) Text(fmt string, arg ...any) { c.Render(render.Text{Format: fmt, Data: arg}) }
func (c *Context) Data(data render.Data)       { c.Render(data) }
func (c *Context) Reader(data render.Reader)   { c.Render(data) }

func (c *Context) Bind(data any) error         { return _Bind{}.Bind(c, data) }
func (c *Context) BindAll(data any) error      { return _BindAll{}.Bind(c, data) }
func (c *Context) BindForm(data any) error     { return _BindAll{}.Bind(c, data) }
func (c *Context) BindPostForm(data any) error { return _BindAll{}.Bind(c, data) }
func (c *Context) BindHeader(data any) error   { return _BindAll{}.Bind(c, data) }
func (c *Context) BindJSON(data any) error     { return _BindAll{}.Bind(c, data) }
func (c *Context) BindURLQuery(data any) error { return _BindAll{}.Bind(c, data) }

func (c *Context) RemoteIP() string    { return c.getRemoteIPs(false)[0] }
func (c *Context) RemoteIPs() []string { return c.getRemoteIPs(false) }

func (c *Context) getRemoteIPs(abort bool) []string {
	var (
		arr       = make([]string, 0, len(c.thz.trustedProxies)+1)
		rIP       = c.fc.RemoteIP()
		isTrusted = false
	)

	for _, v := range c.thz.trustedProxies {
		if v.Contains(rIP) {
			isTrusted = true
			break
		}
	}

	if isTrusted {
		for _, v := range c.thz.trustedHeaders {
			arr = append(arr, pyrokinesis.Bytes.ToString(c.fc.Request.Header.Peek(v)))
			if abort {
				return arr
			}
		}
	}

	return append(arr, rIP.String())
}

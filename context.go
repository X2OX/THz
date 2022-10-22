package THz

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.x2ox.com/THz/render"
	"go.x2ox.com/sorbifolia/httprouter"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type (
	Handler  func(ctx *Context)
	Handlers []Handler
)

type Context struct {
	fc  *fasthttp.RequestCtx
	thz *THz

	mux  sync.RWMutex
	keys map[any]any

	params   httprouter.Params
	handlers Handlers
	index    int
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) IsOptions() bool { return c.fc.Request.Header.IsOptions() }
func (c *Context) IsGet() bool     { return c.fc.Request.Header.IsGet() }
func (c *Context) IsPost() bool    { return c.fc.Request.Header.IsPost() }
func (c *Context) IsHead() bool    { return c.fc.Request.Header.IsHead() }
func (c *Context) IsPut() bool     { return c.fc.Request.Header.IsPut() }
func (c *Context) IsDelete() bool  { return c.fc.Request.Header.IsDelete() }
func (c *Context) IsPatch() bool   { return c.fc.Request.Header.IsPatch() }

func (c *Context) Header(s string) []byte           { return c.fc.Request.Header.Peek(s) }
func (c *Context) SetHeader(k, v string) *Context   { c.fc.Response.Header.Set(k, v); return c }
func (c *Context) SetLocation(v string) *Context    { return c.SetHeader("Location", v) }
func (c *Context) SetContentType(v string) *Context { return c.SetHeader("Content-Type", v) }

func (c *Context) Status(code int) *Context { c.fc.Response.SetStatusCode(code); return c }
func (c *Context) Abort() *Context          { c.index = len(c.handlers); return c }
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

func (c *Context) Param(key string) string { val, _ := c.params.Get(key); return val }

func (c *Context) Bind(data any) error         { return _Bind{}.Bind(c, data) }
func (c *Context) BindAll(data any) error      { return _BindAll{}.Bind(c, data) }
func (c *Context) BindForm(data any) error     { return _BindForm{}.Bind(c, data) }
func (c *Context) BindPostForm(data any) error { return _BindPostForm{}.Bind(c, data) }
func (c *Context) BindHeader(data any) error   { return _BindHeader{}.Bind(c, data) }
func (c *Context) BindJSON(data any) error     { return _BindJSON{}.Bind(c, data) }
func (c *Context) BindURLQuery(data any) error { return _BindURLQuery{}.Bind(c, data) }

func (c *Context) RemoteIP() string    { return c.getRemoteIPs(false)[0] }
func (c *Context) RemoteIPs() []string { return c.getRemoteIPs(false) }

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c == nil || c.fc == nil {
		return
	}
	return c.fc.Deadline()
}
func (c *Context) Done() <-chan struct{} {
	if c == nil || c.fc == nil {
		return nil
	}
	return c.fc.Done()
}
func (c *Context) Err() error {
	if c == nil || c.fc == nil {
		return nil
	}
	return c.fc.Err()
}
func (c *Context) Value(key any) any {
	if c == nil || c.fc == nil {
		return nil
	}
	if val, exists := c.Get(key); exists {
		return val
	}
	return c.fc.Value(key)
}

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

func (c *Context) L() *zap.Logger {
	return c.thz.log.With(
		zap.String(string(c.fc.Method()), c.fc.URI().String()),
	)
}

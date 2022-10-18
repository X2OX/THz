package THz

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.x2ox.com/sorbifolia/httprouter"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type THz struct {
	httprouter.IRouter[Context]

	srv       *fasthttp.Server
	route     *httprouter.Router[Context]
	intercept Handlers

	noRoute Handlers

	trustedProxies []*net.IPNet
	trustedHeaders []string
	ctxPool        sync.Pool

	log *zap.Logger
}

func New() *THz {
	r := httprouter.NewRouter[Context]()
	t := &THz{
		IRouter: r.Group(),
		srv:     &fasthttp.Server{Name: "THz"},
		route:   r,
	}
	t.ctxPool = sync.Pool{New: func() any {
		return &Context{thz: t, index: -1, handlers: t.intercept}
	}}

	t.srv.Handler = t.handle()

	return t
}

func (thz *THz) SetConcurrency(n int) *THz               { thz.srv.Concurrency = n; return thz }
func (thz *THz) SetReadBufferSize(n int) *THz            { thz.srv.ReadBufferSize = n; return thz }
func (thz *THz) SetWriteBufferSize(n int) *THz           { thz.srv.WriteBufferSize = n; return thz }
func (thz *THz) SetReadTimeout(n time.Duration) *THz     { thz.srv.ReadTimeout = n; return thz }
func (thz *THz) SetWriteTimeout(n time.Duration) *THz    { thz.srv.WriteTimeout = n; return thz }
func (thz *THz) SetIdleTimeout(n time.Duration) *THz     { thz.srv.IdleTimeout = n; return thz }
func (thz *THz) SetKeepalivePeriod(n time.Duration) *THz { thz.srv.TCPKeepalivePeriod = n; return thz }
func (thz *THz) SetMaxRequestBodySize(n int) *THz        { thz.srv.MaxRequestBodySize = n; return thz }
func (thz *THz) SetReduceMemoryUsage(n bool) *THz        { thz.srv.ReduceMemoryUsage = n; return thz }

func (thz *THz) SetTrustedProxies(ip ...string) error {
	arr := make([]*net.IPNet, 0, len(ip))
	for _, v := range ip {
		if !strings.Contains(v, "/") {
			_ip := net.ParseIP(v)

			switch len(_ip) {
			case net.IPv4len:
				v += "/32"
			case net.IPv6len:
				v += "/128"
			default:
				return fmt.Errorf("thz: parse ip err: %s", v)
			}
		}

		_, cidr, err := net.ParseCIDR(v)
		if err != nil {
			return err
		}
		arr = append(arr, cidr)
	}

	thz.trustedProxies = arr
	return nil
}

// SetTrustedHeaders
//
// E.g:
// - CF-Connecting-IP
// - X-Forwarded-For
// - X-Real-IP
//
// Remember, order determines priority
func (thz *THz) SetTrustedHeaders(header ...string) {
	thz.trustedHeaders = header
}

func (thz *THz) AddIntercept(intercept ...Handler) {
	// TODO check thz.route
	thz.intercept = append(thz.intercept, intercept...)
}

func (thz *THz) handle() func(c *fasthttp.RequestCtx) {
	return func(c *fasthttp.RequestCtx) {
		ctx := thz.ctxPool.Get().(*Context)
		ctx.fc = c
		ctx.index = -1
		ctx.handlers = append(ctx.handlers, thz.intercept...)
		ctx.keys = make(map[any]any)

		method, uri := httprouter.NewMethod(pyrokinesis.Bytes.ToString(c.Method())),
			pyrokinesis.Bytes.ToString(c.URI().Path())

		handlers := thz.route.Find(method, uri, &ctx.params)
		for _, v := range handlers {
			ctx.handlers = append(ctx.handlers, Handler(v))
		}

		if handlers == nil {
			ctx.handlers = append(ctx.handlers, thz.noRoute...)
		}

		ctx.Next()

		ctx.params = ctx.params[:0]
		ctx.handlers = ctx.handlers[:0]
		ctx.fc = nil
		ctx.keys = nil
		thz.ctxPool.Put(ctx)
	}
}

func (thz *THz) NoRoute(handlers ...Handler) {
	thz.noRoute = handlers
}

func (thz *THz) TestRequest(req *http.Request, msTimeout ...int) (resp *http.Response, err error) {
	timeout := 1000
	if len(msTimeout) > 0 {
		timeout = msTimeout[0]
	}

	if req.Body != http.NoBody && req.Header.Get("Content-Length") == "" {
		req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	}

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}

	dumps := bytes.Split(dump, []byte(" "))
	dumps[1] = []byte(req.URL.String())
	dump = bytes.Join(dumps, []byte(" "))

	conn := new(testConn)

	if _, err = conn.r.Write(dump); err != nil {
		return nil, err
	}

	channel := make(chan error)
	go func() {
		var returned bool
		defer func() {
			if !returned {
				channel <- fmt.Errorf("runtime.Goexit() called in handler or server panic")
			}
		}()

		channel <- thz.srv.ServeConn(conn)
		returned = true
	}()

	if timeout >= 0 {
		select {
		case err = <-channel:
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			return nil, fmt.Errorf("test: timeout error %vms", timeout)
		}
	} else {
		err = <-channel
	}

	if err != nil && err != fasthttp.ErrGetOnly {
		return nil, err
	}

	buffer := bufio.NewReader(&conn.w)

	return http.ReadResponse(buffer, req)
}

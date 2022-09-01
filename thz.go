package THz

import (
	"fmt"
	"net"
	"strings"
	"sync"

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

func (thz *THz) SetLog(log *zap.Logger) { thz.log = log }

func (thz *THz) handle() func(c *fasthttp.RequestCtx) {
	return func(c *fasthttp.RequestCtx) {
		ctx := thz.ctxPool.Get().(*Context)
		ctx.fc = c
		ctx.index = -1
		ctx.handlers = append(ctx.handlers, thz.intercept...)
		ctx.keys = make(map[any]any)

		method, uri := httprouter.NewMethod(pyrokinesis.Bytes.ToString(c.Method())),
			pyrokinesis.Bytes.ToString(c.URI().Path())

		handlers, params := thz.route.Find(method, uri)
		for _, v := range handlers {
			ctx.handlers = append(ctx.handlers, Handler(v))
		}
		ctx.params = params

		ctx.Next()

		ctx.params = ctx.params[:0]
		ctx.handlers = ctx.handlers[:0]
		ctx.fc = nil
		ctx.keys = nil
		thz.ctxPool.Put(ctx)
	}
}

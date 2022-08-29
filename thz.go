package THz

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
	"go.x2ox.com/sorbifolia/httprouter"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type THz struct {
	httprouter.IRouter[Context]

	srv   *fasthttp.Server
	route *httprouter.Router[Context]

	trustedProxies []*net.IPNet
	trustedHeaders []string
	ctxPool        sync.Pool
}

func New() *THz {
	r := httprouter.NewRouter[Context]()
	t := &THz{
		IRouter: r.Group(),
		srv:     &fasthttp.Server{Name: "THz"},
		route:   r,
	}
	t.ctxPool = sync.Pool{New: func() any {
		return &Context{thz: t, index: -1}
	}}

	t.srv.Handler = func(c *fasthttp.RequestCtx) {
		ctx := t.ctxPool.Get().(*Context)
		ctx.fc = c
		ctx.index = -1

		method, uri := httprouter.NewMethod(pyrokinesis.Bytes.ToString(c.Method())),
			pyrokinesis.Bytes.ToString(c.URI().Path())

		ctx.handlers, ctx.params = t.route.Find(method, uri)

		ctx.Next()

		ctx.params = ctx.params[:0]
		ctx.handlers = ctx.handlers[:0]
		ctx.fc = nil
		t.ctxPool.Put(ctx)
	}

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

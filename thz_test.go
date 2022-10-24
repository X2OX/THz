package THz

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
)

func TestNoFound(t *testing.T) {
	thz := New()

	thz.NoRoute(func(c *Context) {
		c.Status(http.StatusNotFound).JSON("404 Not Found")
	})

	thz.GET("/test", func(c *Context) {
		c.JSON("hello world")
	})

	uri := fasthttp.AcquireURI()
	uri.SetPath("/test")

	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetURI(uri)

	thz.TestHandler(&ctx)

	if ctx.Response.StatusCode() != 200 {
		t.Error("exist route is error")
	}

	uri.SetPath("/noRoute")
	ctx.Request.SetURI(uri)
	thz.TestHandler(&ctx)

	if ctx.Response.StatusCode() != 404 {
		t.Error("noRoute is error")
	}
}

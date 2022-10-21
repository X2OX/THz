package cors

import (
	"github.com/valyala/fasthttp"
	"go.x2ox.com/THz"
	"testing"
)

func TestCORSDefault(t *testing.T) {
	thz := THz.New()
	thz.AddIntercept(New().Middleware())

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")

	handler := thz.TestHandler()
	handler(ctx)

	if string(ctx.Response.Header.Peek(headerAllowOrigin)) != "*" {
		t.Error(string(ctx.Response.Header.Peek(headerAllowOrigin)))
	}

	if string(ctx.Response.Header.Peek(headerAllowCredentials)) != "" {
		t.Error(string(ctx.Response.Header.Peek(headerAllowCredentials)))
	}

	if string(ctx.Response.Header.Peek(headerExposeHeaders)) != "" {
		t.Error(string(ctx.Response.Header.Peek(headerExposeHeaders)))
	}

	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("OPTIONS")
	handler(ctx)

	if string(ctx.Response.Header.Peek("Vary")) != "Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods" {
		t.Error(string(ctx.Response.Header.Peek("Vary")))
	}

	if string(ctx.Response.Header.Peek(headerAllowOrigin)) != "*" {
		t.Error(string(ctx.Response.Header.Peek(headerAllowOrigin)))
	}

	if string(ctx.Response.Header.Peek(headerAllowMethods)) != "GET,POST,HEAD,PUT,DELETE,PATCH" {
		t.Error(string(ctx.Response.Header.Peek(headerAllowMethods)))
	}

	if string(ctx.Response.Header.Peek(headerAllowHeaders)) != "Origin,Content-Length,Content-Type,Authorization" {
		t.Error(string(ctx.Response.Header.Peek(headerAllowHeaders)))
	}
}

func TestAllowOrigins(t *testing.T) {
	thz := THz.New()
	thz.AddIntercept(New("1.1.1.1", "192.168.1.1").Middleware())

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.Set(headerOrigin, "10.10.1.1")
	ctx.Request.Header.SetMethod("GET")

	handler := thz.TestHandler()
	handler(ctx)

	if string(ctx.Response.Header.Peek("Vary")) != "Origin" {
		t.Error(string(ctx.Response.Header.Peek("Vary")))
	}

	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.Set(headerOrigin, "1.1.1.1")
	ctx.Request.Header.SetMethod("OPTIONS")

	handler(ctx)

	if string(ctx.Response.Header.Peek("Vary")) != "Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods" {
		t.Error(string(ctx.Response.Header.Peek("Vary")))
	}
}

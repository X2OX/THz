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

	thz.TestHandler(ctx)

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
	thz.TestHandler(ctx)

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
	thz.AddIntercept(New("http:", "https:").Middleware())

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.Set(headerOrigin, "chrome-extension://")
	ctx.Request.Header.SetMethod("GET")

	thz.TestHandler(ctx)

	if string(ctx.Response.Header.Peek("Vary")) != "Origin" {
		t.Error(string(ctx.Response.Header.Peek("Vary")))
	}

	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.Set(headerOrigin, "http:")
	ctx.Request.Header.SetMethod("OPTIONS")

	thz.TestHandler(ctx)

	if string(ctx.Response.Header.Peek("Vary")) != "Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods" {
		t.Error(string(ctx.Response.Header.Peek("Vary")))
	}
}

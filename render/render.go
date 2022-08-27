package render

import (
	"github.com/valyala/fasthttp"
)

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on.
type Render interface {
	// Render writes data with custom ContentType.
	Render(w *fasthttp.RequestCtx) error
	// WriteContentType writes custom ContentType.
	WriteContentType(w *fasthttp.RequestCtx)
}

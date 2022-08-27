package render

import (
	"github.com/valyala/fasthttp"
)

// Data contains ContentType and bytes data.
type Data struct {
	ContentType string
	Data        []byte
}

// Render (Data) writes data with custom ContentType.
func (r Data) Render(w *fasthttp.RequestCtx) (err error) {
	w.SetBody(r.Data)
	return
}

// WriteContentType (Data) writes custom ContentType.
func (r Data) WriteContentType(w *fasthttp.RequestCtx) {
	w.SetContentType(r.ContentType)
}

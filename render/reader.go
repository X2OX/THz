package render

import (
	"io"

	"github.com/valyala/fasthttp"
)

// Reader contains the IO reader and its length, and custom ContentType and other headers.
type Reader struct {
	ContentType   string
	ContentLength int
	Reader        io.Reader
	Headers       map[string]string
}

// Render (Reader) writes data with custom ContentType and headers.
func (r Reader) Render(w *fasthttp.RequestCtx) (err error) {
	for k, v := range r.Headers {
		w.Response.Header.Set(k, v)
	}

	w.Response.SetBodyStream(r.Reader, r.ContentLength)
	return
}

// WriteContentType (Reader) writes custom ContentType.
func (r Reader) WriteContentType(w *fasthttp.RequestCtx) {
	w.SetContentType(r.ContentType)
}

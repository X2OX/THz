package render

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// Text contains the given interface object slice and its format.
type Text struct {
	Format string
	Data   []any
}

const plainContentType = "text/plain; charset=utf-8"

// Render (Text) writes data with custom ContentType.
func (r Text) Render(w *fasthttp.RequestCtx) error {
	if len(r.Data) > 0 {
		_, err := fmt.Fprintf(w, r.Format, r.Data...)
		return err
	}
	w.SetBodyString(r.Format)

	return nil
}

// WriteContentType (Text) writes Plain ContentType.
func (r Text) WriteContentType(w *fasthttp.RequestCtx) {
	w.SetContentType(plainContentType)
}

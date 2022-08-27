package render

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type JSON struct {
	Data any
}

const jsonContentType = "application/json; charset=utf-8"

// Render (JSON) writes data with custom ContentType.
func (r JSON) Render(w *fasthttp.RequestCtx) error {
	data, err := json.Marshal(r.Data)
	if err == nil {
		w.SetBody(data)
	}
	return err
}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSON) WriteContentType(w *fasthttp.RequestCtx) {
	w.SetContentType(jsonContentType)
}

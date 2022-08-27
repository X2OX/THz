package header

import (
	"testing"

	"github.com/valyala/fasthttp"
)

type SDA struct {
	A string  `header:"Content-Type"`
	B string  `header:"b"`
	C float64 `header:"c"`
	D *int    `header:"d"`
}

func TestParse(t *testing.T) {
	h := &fasthttp.RequestHeader{}
	h.Set("Content-Type", "json")
	h.Set("b", "1")
	h.Set("c", "1.1")
	h.Set("d", "1")

	vvv := &SDA{}
	if err := Parse(h, vvv); err != nil {
		t.Error(err)
	}
}

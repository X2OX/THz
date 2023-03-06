package header

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

type SDA struct {
	A string     `header:"Content-Type"`
	B string     `header:"b"`
	C float64    `header:"c"`
	D *int       `header:"d"`
	T *time.Time `header:"T"`
}

func TestParse(t *testing.T) {
	h := &fasthttp.RequestHeader{}
	h.Set("Content-Type", "json")
	h.Set("b", "1")
	h.Set("c", "1.1")
	h.Set("d", "1")
	h.Set("t", time.Now().Format(time.RFC1123))

	vvv := &SDA{}
	if err := Parse(h, vvv); err != nil {
		t.Error(err)
	}
}

func TestParseMap(t *testing.T) {
	type tm struct {
		B map[string]string `header:"b"`
	}

	h := &fasthttp.RequestHeader{}
	h.Set("b", `{"a":"123","q":"456"}`)

	vvv := &tm{}
	if err := Parse(h, vvv); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual("123", vvv.B["a"]) || !reflect.DeepEqual("456", vvv.B["q"]) {
		t.Error("fail to parse map")
	}
}

func TestParseSlice(t *testing.T) {
	type ta struct {
		B [3]string `header:"b"`
	}

	type ts struct {
		B []string `header:"b"`
	}

	h := &fasthttp.RequestHeader{}
	h.Set("b", `["123","456","789"]`)

	tSlice := &ts{}
	if err := Parse(h, tSlice); err != nil {
		t.Error(err)
	}

	res := []string{"123", "456", "789"}
	for i, v := range tSlice.B {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if !reflect.DeepEqual(v, res[i]) {
				t.Errorf("expect: %v,get: %v", v, res[i])
			}
		})
	}

	tArr := &ta{}
	if err := Parse(h, tArr); err != nil {
		t.Error(err)
	}
	for i, v := range tArr.B {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if !reflect.DeepEqual(v, res[i]) {
				t.Errorf("expect: %v,get: %v", v, res[i])
			}
		})
	}
}

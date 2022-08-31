package main

import (
	"fmt"

	"go.x2ox.com/THz"
)

var h = func(ctx *THz.Context) {
	ctx.Set("asd", "asd")
	ctx.JSON(map[string]string{"hello": "string"})
}

func main() {
	t := THz.New()
	t.AddIntercept(func(ctx *THz.Context) {
		ctx.SetHeader("123", "asd")
	})
	t.Any("/", h)
	if err := t.ListenAndServe(":8888"); err != nil {
		fmt.Println(err)
	}

}

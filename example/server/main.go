package main

import (
	"fmt"

	"go.x2ox.com/THz"
	"go.x2ox.com/THz/middleware/cors"
	"go.x2ox.com/THz/middleware/recovery"
)

var h = func(ctx *THz.Context) {
	ctx.Set("asd", "asd")
	ctx.JSON(map[string]string{"hello": "string"})
}

func main() {
	t := THz.New()
	t.AddIntercept(recovery.New().Middleware())
	t.AddIntercept(cors.New().Middleware())
	t.AddIntercept(func(ctx *THz.Context) {
		ctx.SetHeader("123", "asd")
	})
	t.Any("/", h)
	t.Any("/1", func(c *THz.Context) {
		panic("test")
	})
	if err := t.ListenAndServe(":8888"); err != nil {
		fmt.Println(err)
	}

}

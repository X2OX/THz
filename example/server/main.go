package main

import (
	"fmt"

	"go.x2ox.com/THz"
)

var h = func(ctx *THz.Context) {
	ctx.JSON(map[string]string{"hello": "string"})
}

func main() {
	t := THz.New()

	t.Any("/", h)
	if err := t.ListenAndServe(":8888"); err != nil {
		fmt.Println(err)
	}

}

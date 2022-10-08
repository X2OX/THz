package THz

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {
	thz := New()
	thz.GET("/test", func(t *Context) {
		t.JSON("hello world")
	})

	go func(t *THz) {
		time.Sleep(5 * time.Second)
		if err := t.Stop(); err != nil {
			fmt.Println(err)
		}
	}(thz)

	if err := thz.ListenAndServe(":8080"); err != nil {
		fmt.Println("start err")
	}
}

func TestRunAndStop(t *testing.T) {
	thz := New()
	thz.GET("/run", func(t *Context) {
		t.JSON("worked")
	})

	go func(t *THz) {
		time.Sleep(5 * time.Second)
		if err2 := t.Stop(); err2 != nil {
			fmt.Println("stop err")
			fmt.Println(err2)
		}
	}(thz)

	ln, err := net.Listen("tcp4", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = thz.Run(ln); err != nil {
		fmt.Println(err)
	}

}

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

	go func() {
		time.Sleep(5 * time.Second)
		if err := thz.Stop(); err != nil {
			fmt.Println(err)
		}
	}()

	if err := thz.ListenAndServe(":8080"); err != nil {
		fmt.Println("start err")
	}
}

func TestRunAndStop(t *testing.T) {
	thz := New()
	thz.GET("/run", func(t *Context) {
		t.JSON("worked")
	})

	go func() {
		time.Sleep(1 * time.Minute)
		if err2 := thz.Stop(); err2 != nil {
			fmt.Println("stop err")
			fmt.Println(err2)
		}
	}()

	ln, err := net.Listen("tcp4", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = thz.Run(ln); err != nil {
		fmt.Println(err)
	}
}

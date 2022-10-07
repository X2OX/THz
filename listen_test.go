package THz

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {
	thz := New()
	thz.GET("/test", func(t *Context) {
		t.JSON("hello world")
	})

	if err := thz.ListenAndServe(":8080"); err != nil {
		fmt.Println("start err")
	}
}

func TestRunAndStop(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
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
		wg.Done()
	}()

	ln, err := net.Listen("tcp4", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = thz.Run(ln); err != nil {
		fmt.Println(err)
	}

	wg.Wait()
	fmt.Println("finished")
}

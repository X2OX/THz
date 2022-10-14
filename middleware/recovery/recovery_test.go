package recovery

import (
	"fmt"
	"go.x2ox.com/THz"
	"net/http"
	"testing"
	"time"
)

func TestRecovery(t *testing.T) {
	thz := THz.New()
	thz.AddIntercept(New().Middleware())
	thz.GET("/recovery", func(_ *THz.Context) {
		panic("test recovery")
	})

	ch := make(chan int)
	go func() {
		<-ch
		if err := thz.Stop(); err != nil {
			fmt.Println("error stop")
		}
	}()

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 3; i++ {
			resp, err := http.Get("http://localhost:8081/recovery")
			if err != nil {
				fmt.Println("get error", err)
				return
			}

			if resp != nil {
				fmt.Println(resp)
			}
		}

		ch <- 1
	}()

	if err := thz.ListenAndServe(":8081"); err != nil {
		fmt.Println("start err")
	}
}

package recovery

import (
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
			t.Error(err)
		}
	}()

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 3; i++ {
			resp, err := http.Get("http://localhost:8081/recovery")
			if err != nil {
				t.Error(err)
				return
			}

			if resp != nil {
				t.Log(resp)
			}
		}

		ch <- 1
	}()

	if err := thz.ListenAndServe(":8081"); err != nil {
		t.Error(err)
	}
}

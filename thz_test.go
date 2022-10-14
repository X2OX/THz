package THz

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"io"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestHandleNotFound(t *testing.T) {
	thz := New()
	thz.GET("/test", func(c *Context) {
		c.JSON("hello world")
	})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		wg.Wait()
		if err := thz.Stop(); err != nil {
			t.Error(err)
		}
	}()

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 10; i++ {
			resp, err := http.Get("http://localhost:8080/test")
			if err != nil {
				t.Error(err)
			}

			if resp != nil {
				var buf buffer.Buffer
				if _, err = io.Copy(&buf, resp.Body); err != nil {
					t.Error(err)
				}

				t.Log(fmt.Sprintf("%s--%s", strconv.Itoa(resp.StatusCode), string(buf.Bytes())))
			}
		}

		wg.Done()
	}()

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 10; i++ {
			resp, err := http.Get("http://localhost:8080/te")
			if err != nil {
				t.Error(err)
			}

			if resp != nil {
				t.Log(resp.StatusCode)
			}
		}

		wg.Done()
	}()

	if err := thz.ListenAndServe(":8080"); err != nil {
		t.Error(err)
	}
}

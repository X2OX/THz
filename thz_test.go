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

func TestNoFound(t *testing.T) {
	thz := New()

	thz.NoRoute(func(c *Context) {
		c.Status(http.StatusNotFound).JSON("404 Not Found")
	})

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

			if resp != nil && resp.StatusCode != http.StatusOK {
				var buf buffer.Buffer
				if _, err = io.Copy(&buf, resp.Body); err != nil {
					t.Error(err)
				}

				t.Logf(fmt.Sprintf("%s--%s", strconv.Itoa(resp.StatusCode), string(buf.Bytes())))
				t.Error("There are some errors")
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

			if resp != nil && resp.StatusCode != http.StatusNotFound {
				t.Log(resp.StatusCode)
				t.Error("There are some errors")
			}
		}

		wg.Done()
	}()

	if err := thz.ListenAndServe(":8080"); err != nil {
		t.Error(err)
	}
}

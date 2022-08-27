package THz

import (
	"fmt"
	"net/url"
	"testing"
)

func TestA(t *testing.T) {
	u, err := url.Parse("/asdasd/asdasdd/asasd?aaa=ava&asdas=123")
	fmt.Println(err)
	fmt.Println(u)

	// (&Context{}).fc.Request.Header.Peek()
}

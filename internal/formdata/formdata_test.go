package formdata

import (
	"bytes"
	"mime/multipart"
	"testing"
	"time"
)

type SDA struct {
	A string                  `form:"a"`
	B string                  `form:"b"`
	C string                  `form:"c"`
	D []*multipart.FileHeader `form:"d"`
	T *time.Time              `form:"t"`
}

func TestParse(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	{
		_ = mw.SetBoundary("go.x2ox.com/THz")
		_ = mw.WriteField("a", "123")
		_ = mw.WriteField("b", "123")
		_ = mw.WriteField("c", "123")
		_ = mw.WriteField("c", "123")
		_ = mw.WriteField("t", time.Now().Format(time.RFC3339))

		iw, _ := mw.CreateFormFile("d", "111.txt")
		_, _ = iw.Write([]byte("123"))
		iw, _ = mw.CreateFormFile("d", "222.txt")
		_, _ = iw.Write([]byte("123"))

		_ = mw.Close()
	}

	mr := multipart.NewReader(buf, "go.x2ox.com/THz")
	mf, err := mr.ReadForm(10 << 20)
	if err != nil {
		t.Error(err)
		return
	}
	vvv := &SDA{}
	if err = Parse(mf, vvv); err != nil {
		t.Error(err)
	}
}

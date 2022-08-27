package THz

import (
	"encoding/json"
	"errors"
	"net/url"

	"go.x2ox.com/THz/internal/formdata"
	"go.x2ox.com/THz/internal/header"
	"go.x2ox.com/THz/internal/urlquery"
	"go.x2ox.com/sorbifolia/pyrokinesis"
)

type Binding interface {
	Bind(*Context, any) error
}

type (
	_Bind    struct{}
	_BindAll struct{}

	_BindForm     struct{}
	_BindPostForm struct{}
	_BindURLQuery struct{}
	_BindHeader   struct{}
	_BindJSON     struct{}
)

func (_Bind) Bind(c *Context, a any) error {
	contentType := c.fc.Request.Header.ContentType()
	if len(contentType) == 0 {
		return _BindURLQuery{}.Bind(c, a)
	}

	switch filterFlags(pyrokinesis.Bytes.ToString(contentType)) {
	case "application/json":
		return _BindJSON{}.Bind(c, a)
	case "application/x-www-form-urlencoded":
		return _BindForm{}.Bind(c, a)
	case "multipart/form-data":
		return _BindPostForm{}.Bind(c, a)
	}

	return errors.New("ContentType no support")
}

func (_BindAll) Bind(c *Context, a any) error {
	if err := (_Bind{}).Bind(c, a); err != nil {
		return err
	}
	if err := (_BindURLQuery{}).Bind(c, a); err != nil {
		return err
	}
	if err := (_BindHeader{}).Bind(c, a); err != nil {
		return err
	}
	return nil
}

func (_BindJSON) Bind(c *Context, a any) error {
	body, err := c.fc.Request.BodyUncompressed()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, a)
}

func (_BindHeader) Bind(c *Context, a any) error { return header.Parse(&c.fc.Request.Header, a) }

func (_BindForm) Bind(c *Context, a any) error {
	f, err := c.fc.Request.MultipartForm()
	if err != nil {
		return err
	}
	return formdata.Parse(f, a)
}

func (_BindURLQuery) Bind(c *Context, a any) error {
	u, err := url.Parse(pyrokinesis.Bytes.ToString(c.fc.Request.Header.RequestURI()))
	if err != nil {
		return err
	}

	var val url.Values
	if val, err = url.ParseQuery(u.RawQuery); err != nil {
		return err
	}

	return urlquery.Parse(val, a)
}

func (_BindPostForm) Bind(c *Context, a any) error {
	u, err := url.Parse(pyrokinesis.Bytes.ToString(c.fc.Request.Body()))
	if err != nil {
		return err
	}

	var val url.Values
	if val, err = url.ParseQuery(u.RawQuery); err != nil {
		return err
	}

	return urlquery.Parse(val, a)
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

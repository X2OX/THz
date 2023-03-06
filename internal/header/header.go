package header

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
	"go.x2ox.com/sorbifolia/pyrokinesis"
	"go.x2ox.com/sorbifolia/strong"
)

func Parse(data *fasthttp.RequestHeader, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("unexpected")
	}

	return parse(data, rv)
}

func parse(data *fasthttp.RequestHeader, v reflect.Value) error {
	if data.Len() == 0 {
		return nil
	}

	switch v.Kind() {
	case reflect.Interface:
		return parse(data, v.Elem())
	// case reflect.Map: // TODO: support map
	case reflect.Pointer:
		if v.IsNil() {
			if !v.CanSet() {
				return nil
			}
			v.Set(reflect.New(v.Type().Elem()))
		}
		return parse(data, v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			t := v.Type().Field(i)

			if t.Anonymous && v.Field(i).Kind() == reflect.Struct { // Anonymous Struct
				if err := parse(data, v.Field(i)); err != nil {
					return err
				}
				continue
			}

			tag := t.Tag.Get("header")
			if tag == "" {
				continue
			}

			if err := setValue(v.Field(i), pyrokinesis.Bytes.ToString(data.Peek(tag))); err != nil {
				return err
			}
		}
	default:
		return errors.New("unknown field type")
	}
	return nil
}

func setValue(v reflect.Value, data string) error {
	if len(data) == 0 || !v.CanAddr() {
		return nil
	}
	switch v.Kind() {
	case reflect.Bool:
		val, err := strong.Parse[bool](data)
		if err != nil {
			return err
		}
		v.SetBool(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strong.Parse[int64](data)
		if err != nil {
			return err
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strong.Parse[uint64](data)
		if err != nil {
			return err
		}
		v.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return err
		}
		v.SetFloat(val)
	case reflect.Complex64, reflect.Complex128:
		val, err := strconv.ParseComplex(data, 128)
		if err != nil {
			return err
		}
		v.SetComplex(val)
	case reflect.Interface:
		v.Set(reflect.ValueOf(data))
	case reflect.Pointer:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return setValue(v.Elem(), data)
	case reflect.Array:
		return json.Unmarshal([]byte(data), v.Addr().Interface())
	case reflect.Slice:
		return json.Unmarshal([]byte(data), v.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal([]byte(data), v.Addr().Interface())
	case reflect.String:
		v.SetString(data)
	case reflect.Struct:
		switch v.Type() {
		case reflect.TypeOf(time.Time{}):
			if tn, _ := strong.Parse[int64](data); tn != 0 {
				v.Set(reflect.ValueOf(time.Unix(tn, 0)))
				return nil
			}

			if t, err := time.Parse(time.RFC3339, data); err == nil {
				v.Set(reflect.ValueOf(t))
				return nil
			}

			if t, err := time.Parse(time.RFC1123, data); err == nil {
				v.Set(reflect.ValueOf(t))
				return nil
			}

			return errors.New("parse time error")
		default:
			return errors.New("unknown field type")
		}
	default:
		return errors.New("unknown field type")
	}
	return nil
}

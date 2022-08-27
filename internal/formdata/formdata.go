package formdata

import (
	"errors"
	"mime/multipart"
	"reflect"
	"strconv"

	"go.x2ox.com/sorbifolia/strong"
)

var (
	mfType = reflect.TypeOf(&multipart.FileHeader{})
)

func Parse(data *multipart.Form, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("unexpected")
	}

	return parse(data, rv)
}

func parse(data *multipart.Form, v reflect.Value) error {
	if len(data.File) == 0 && len(data.Value) == 0 {
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

			tag := t.Tag.Get("form")
			if tag == "" {
				continue
			}
			if val, ok := data.File[tag]; ok {
				if err := setFileValue(v.Field(i), val); err != nil {
					return err
				}
			}
			if val, ok := data.Value[tag]; ok {
				if err := setValue(v.Field(i), val); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("unknown field type")
	}
	return nil
}

func setValue(v reflect.Value, data []string) error {
	if len(data) == 0 || !v.CanAddr() {
		return nil
	}
	switch v.Kind() {
	case reflect.Bool:
		val, err := strong.Parse[bool](data[0])
		if err != nil {
			return err
		}
		v.SetBool(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strong.Parse[int64](data[0])
		if err != nil {
			return err
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strong.Parse[uint64](data[0])
		if err != nil {
			return err
		}
		v.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			return err
		}
		v.SetFloat(val)
	case reflect.Complex64, reflect.Complex128:
		val, err := strconv.ParseComplex(data[0], 128)
		if err != nil {
			return err
		}
		v.SetComplex(val)
	case reflect.Array:
		for i, l := 0, v.Len(); i < l; i++ {
			if len(data) < i+1 {
				return nil
			}
			if err := setValue(v.Index(i), data[i:]); err != nil {
				return err
			}
		}
	case reflect.Interface:
		if len(data) == 1 {
			v.Set(reflect.ValueOf(data[0]))
			return nil
		}
		v.Set(reflect.ValueOf(data))
	case reflect.Pointer:
		if v.IsNil() {
			if !v.CanSet() {
				return nil
			}
			v.Set(reflect.New(v.Type().Elem()))
		}
		return setValue(v.Elem(), data)
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), len(data), len(data)))
		for i, l := 0, v.Len(); i < l; i++ {
			if err := setValue(v.Index(i), data[i:]); err != nil {
				return err
			}
		}
	case reflect.String:
		v.SetString(data[0])
	default:
		return errors.New("unknown field type")
	}
	return nil
}

func setFileValue(v reflect.Value, data []*multipart.FileHeader) error {
	if len(data) == 0 || !v.CanAddr() {
		return nil
	}
	switch v.Kind() {
	case reflect.Array:
		for i, l := 0, v.Len(); i < l; i++ {
			if len(data) < i+1 {
				return nil
			}
			if err := setFileValue(v.Index(i), data[i:]); err != nil {
				return err
			}
		}
	case reflect.Interface:
		if len(data) == 1 {
			v.Set(reflect.ValueOf(data[0]))
			return nil
		}
		v.Set(reflect.ValueOf(data))
	case reflect.Pointer:
		if v.Type() != mfType {
			return errors.New("unknown type")
		}
		v.Set(reflect.ValueOf(data[0]))
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), len(data), len(data)))
		for i, l := 0, v.Len(); i < l; i++ {
			if err := setFileValue(v.Index(i), data[i:]); err != nil {
				return err
			}
		}
	default:
		return errors.New("unknown field type")
	}
	return nil
}

package model

import (
	"reflect"
)

type LocalString map[string]string

func ToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if tagv := f.Tag.Get("json"); tagv != "" && tagv != "-" {
			result[tagv] = v.Field(i).Interface()
		}
	}

	return result
}

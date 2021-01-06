package mysql2

import (
	"errors"
	"reflect"
)

//  判断struct是单数还是复数，只有两种格式是对的*struct,  *[]*struct
func isBmMany(data interface{}) (bool, error) {
	v := reflect.TypeOf(data)
	return _isBmMany(v)
}

func _isBmMany(v reflect.Type) (bool, error) {
	typeErrRead := "struct bm只有两种形式是对的，*struct,  *[]*struct"
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	} else {
		return false, errors.New(typeErrRead)
	}
	switch v.Kind() {
	case reflect.Slice:
		res, err := _isBmMany(v.Elem())
		if err != nil {
			return false, err
		}
		if !res {
			return true, nil
		}
		return false, errors.New(typeErrRead)

	case reflect.Struct:
		return false, nil

	}
	return false, errors.New(typeErrRead)
}

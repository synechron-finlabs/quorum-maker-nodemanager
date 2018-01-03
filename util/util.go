package util

import (
	"reflect"
	"strings"
	"strconv"
	"fmt"
)

func TakeSliceArg(arg interface{}) (out []interface{}, ok bool) {
	slice, success := takeArg(arg, reflect.Slice)
	if !success {
		ok = false
		return
	}
	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

func HexStringtoInt64(hexval string) (intval int64) {
	hexval = strings.TrimSuffix(hexval, "\n")
	hexval = strings.TrimPrefix(hexval, "0x")
	intval, err := strconv.ParseInt(hexval, 16, 64)
	if err != nil {
		fmt.Println(err)
	}
	return intval
}

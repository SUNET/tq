package utils

import (
	"reflect"
	"runtime"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Keys(x map[string]interface{}) []string {
	keys := make([]string, 0, len(x))
	for k, _ := range x {
		keys = append(keys, k)
	}
	return keys
}

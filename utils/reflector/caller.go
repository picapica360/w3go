package reflector

import (
	"reflect"
	"runtime"
)

// GetFuncName get function name
// i -> is the function.
func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

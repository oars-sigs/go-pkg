package flow

import (
	"reflect"

	"github.com/antonmedv/expr"
)

func Eval(s string, vars *Gvars) (interface{}, error) {
	envs := vars.Vars()
	envs["getOne"] = getOne
	vars.mutex.Lock()
	defer vars.mutex.Unlock()
	return expr.Eval(s, envs)
}

func getOne(p interface{}) interface{} {
	v := reflect.ValueOf(p)
	for _, k := range v.MapKeys() {
		return v.MapIndex(k).Interface()
	}
	return ""
}

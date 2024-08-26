package flow

import (
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

func Eval(s string, vars *Gvars) (interface{}, error) {
	envs := vars.Vars()
	envs["getOne"] = getOne
	vars.mutex.Lock()
	defer vars.mutex.Unlock()
	p, err := expr.Compile(s, expr.Env(envs))
	if err != nil {
		return nil, err
	}
	return vm.Run(p, envs)
}

func getOne(p interface{}) interface{} {
	v := reflect.ValueOf(p)
	for _, k := range v.MapKeys() {
		return v.MapIndex(k).Interface()
	}
	return ""
}

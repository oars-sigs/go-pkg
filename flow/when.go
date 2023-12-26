package flow

import "fmt"

func when(s string, vars *Gvars) bool {
	res, err := Eval(s, vars)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if v, ok := res.(bool); ok {
		return v
	}
	return false
}

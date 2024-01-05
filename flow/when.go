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

type whenAction struct{}

func (a *whenAction) Do(conf *Config, params interface{}) (interface{}, error) {
	return "slip", nil
}

func (a *whenAction) Params() interface{} {
	return nil
}

func (a *whenAction) Scheme() string {
	return ""
}

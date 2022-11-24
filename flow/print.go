package flow

import "fmt"

type PrintAction string

func (a *PrintAction) Do(conf *Config, params interface{}) (interface{}, error) {
	fmt.Println(params)
	return nil, nil
}
func (a *PrintAction) Params() interface{} {
	return nil
}

func (a *PrintAction) Scheme() string {
	return ""
}

package flow

import (
	"fmt"
	"io/ioutil"

	"github.com/dop251/goja"
)

var JsFunMap = map[string]map[string]interface{}{
	"console": {
		"log": func(v ...any) {
			fmt.Print("console.log: ")
			fmt.Println(v...)
		},
	},
}

func JSSet(n, k string, v interface{}) {
	if _, ok := JsFunMap[n]; !ok {
		JsFunMap[n] = make(map[string]interface{})
	}
	JsFunMap[n][k] = v
}

type JsAction struct {
	Script string      `yaml:"script"`
	Input  interface{} `yaml:"args"`
	File   string      `yaml:"file"`
}

func (a *JsAction) Do(conf *Config, params interface{}) (interface{}, error) {
	args := params.(JsAction)
	jsvm := goja.New()
	var output interface{}
	JSSet("sys", "args", args.Input)
	JSSet("sys", "output", func(v any) {
		output = v
	})
	for k, v := range JsFunMap {
		jsvm.Set(k, v)
	}
	if args.File != "" {
		data, err := ioutil.ReadFile(args.File)
		if err != nil {
			return nil, err
		}
		args.Script = string(data)
	}
	_, err := jsvm.RunString(args.Script)
	if err != nil {
		return nil, err
	}
	return output, nil
}
func (a *JsAction) Params() interface{} {
	return JsAction{}
}

func (a *JsAction) Scheme() string {
	return ""
}

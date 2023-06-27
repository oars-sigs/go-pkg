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
	//fns :=make([]string,0)
	var output interface{}
	JSSet("sys", "args", args.Input)
	JSSet("sys", "output", func(v any) {
		output = v
	})
	var fns []func()
	JSSet("sys", "defer", func(fn func()) {
		fns = append(fns, fn)
	})
	JSSet("sys", "gdefer", func(fn func()) {
		conf.PTasks = append(conf.PTasks, fn)
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
	for i := len(fns); i > 0; i-- {
		fns[i-1]()
	}
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

package flow

import (
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {
	task := Task{
		"tasks": []Task{
			Task{
				"print": "hello world",
			},
			Task{
				"print": "{{ .ctx.item }}{{ .ctx.itemKey }}",
			},
		},
		"loop": []string{"a", "b"},
	}
	// task := Task{
	// 	"print": "ssss{{ .ctx.item }}",
	// 	"loop":  []string{"a", "b"},
	// }
	printAct := PrintAction("")
	AddCustomActions("print", &printAct)
	awaits := newAwait()
	c, err := task.Action(nil, awaits, newGvars(&Vars{}))
	if err != nil {
		t.Fatal(err)
		return
	}
	c.Do()
}

type PrintAction string

func (a *PrintAction) Do(conf *Config, params interface{}) (interface{}, error) {
	//args := params.(string)
	fmt.Println(params)
	return nil, nil
}
func (a *PrintAction) Params() interface{} {
	return ""
}

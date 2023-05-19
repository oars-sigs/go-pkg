package flow

import (
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
	c, err := task.Action(&Config{}, newAwait(), NewGvars(&Vars{}))
	if err != nil {
		t.Fatal(err)
		return
	}
	c.Do()
}

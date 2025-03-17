package flow

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func Import(path string, value map[string]any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var p Playbook
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	conf := &Config{
		Next:    p.Next,
		Workdir: ".",
	}
	fmt.Println(value)
	fmt.Println(p.Values)
	p.Values = MergeValues(p.Values, value)
	fmt.Println(p.Values)
	vars := &Vars{
		Values: p.Values,
		Ctx:    make(map[string]interface{}),
	}
	p.gvars = NewGvars(vars)
	return p.Run(conf)
}

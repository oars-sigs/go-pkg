package flow

import (
	"os"

	"gopkg.in/yaml.v2"
)

func Import(path string, value map[string]any) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return value, err
	}
	var p Playbook
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return value, err
	}
	conf := &Config{
		Next:    p.Next,
		Workdir: ".",
	}
	p.Values = MergeValues(p.Values, value)
	vars := &Vars{
		Values: p.Values,
		Ctx:    make(map[string]interface{}),
	}
	p.gvars = NewGvars(vars)
	err = p.Run(conf)
	if err != nil {
		return value, err
	}
	return p.Values, nil
}

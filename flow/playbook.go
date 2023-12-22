package flow

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func Run(path string, valuePath ...string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var p Playbook
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	p.await = newAwait()
	conf := &Config{
		Next:    p.Next,
		Workdir: ".",
	}
	if p.Values == nil {
		p.Values = make(map[string]interface{})
	}
	if len(valuePath) > 0 && valuePath[0] != "" {
		data, err := os.ReadFile(valuePath[0])
		if err != nil {
			return err
		}
		var value map[string]interface{}
		err = yaml.Unmarshal(data, &value)
		if err != nil {
			return err
		}
		for k, v := range value {
			p.Values[k] = v
		}
	}
	vars := &Vars{
		Values: p.Values,
		Ctx:    make(map[string]interface{}),
	}
	p.gvars = NewGvars(vars)
	for _, m := range p.Modules {
		AddCustomActions(m.Name, &m)
	}
	AddCustomActions("print", new(PrintAction))
	return p.Run(conf)
}

type Playbook struct {
	Tasks     []Task                 `yaml:"tasks"`
	Values    map[string]interface{} `yaml:"values"`
	Modules   []Module               `yaml:"modules"`
	gvars     *Gvars
	index     int
	await     *gawait
	deferTask []customAction
}

func NewPlaybook(tasks []Task, vars *Gvars) *Playbook {
	return &Playbook{
		Tasks: tasks,
		gvars: vars,
		await: newAwait(),
	}
}

func (p *Playbook) Run(conf *Config) error {
	for {
		if p.index == len(p.Tasks) {
			break
		}
		t := p.Tasks[p.index]
		c, err := t.Action(conf, p.await, p.gvars)
		if err != nil {
			return err
		}
		if len(c.DeferTasks) > 0 {
			defer func() {
				_, err = c.Do()
				if err != nil {
					fmt.Println(err)
				}
			}()
		} else {
			_, err = c.Do()
			if err != nil {
				return err
			}
		}

		p.index = p.index + 1
	}
	for i := len(conf.PTasks); i > 0; i-- {
		conf.PTasks[i-1]()
	}
	return nil

}

func (p *Playbook) Next(id string, conf *Config, vars *Gvars) (interface{}, error) {
	task, index, err := p.getTask(id)
	if err != nil {
		return nil, err
	}
	p.index = index
	c, err := task.Action(conf, p.await, vars)
	if err != nil {
		return nil, err
	}
	return c.Do()
}

func (p *Playbook) getTask(id string) (Task, int, error) {
	for index, t := range p.Tasks {
		if tid, ok := t["id"]; ok {
			if tid == id {
				return t, index, nil
			}
		}
	}
	return nil, 0, errors.New("task not found")
}

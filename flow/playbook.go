package flow

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Run(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var p Playbook
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	conf := &Config{
		Next: p.Next,
	}
	if p.Values == nil {
		p.Values = make(map[string]interface{})
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
	Tasks   []Task                 `yaml:"tasks"`
	Values  map[string]interface{} `yaml:"values"`
	Modules []Module               `yaml:"modules"`
	gvars   *Gvars
	index   int
	await   *gawait
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
		_, err = c.Do()
		if err != nil {
			return err
		}
		p.index = p.index + 1
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

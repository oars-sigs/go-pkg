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
		Next: p.next,
	}
	return p.Run(conf, newAwait(), newGvars(&p.Vars))
}

type Playbook struct {
	Tasks []Task `yaml:"tasks"`
	Vars  Vars   `yaml:"vars"`
	index int
}

func (p *Playbook) Run(conf *Config, await *gawait, vars *gvars) error {
	// conf := &Config{
	// 	Next: p.next,
	// }
	// await := newAwait()
	// vars := newGvars(&Vars{})
	for {
		if p.index == len(p.Tasks) {
			break
		}
		t := p.Tasks[p.index]
		c, err := t.Action(conf, await, vars)
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

func (p *Playbook) next(id string, conf *Config, await *gawait, vars *gvars) (interface{}, error) {
	task, index, err := p.getTask(id)
	if err != nil {
		return nil, err
	}
	p.index = index
	c, err := task.Action(conf, await, vars)
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

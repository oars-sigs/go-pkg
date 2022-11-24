package flow

type Module struct {
	Name   string      `yaml:"name"`
	Tasks  []Task      `yaml:"tasks"`
	Output interface{} `yaml:"output"`
}

func (a *Module) Do(conf *Config, params interface{}) (interface{}, error) {
	vars := &Vars{
		Values: params.(map[string]interface{}),
		Ctx:    make(map[string]interface{}),
	}
	gv := NewGvars(vars)
	playbook := NewPlaybook(a.Tasks, gv)
	config := &Config{
		Workdir: conf.Workdir,
		Next:    playbook.Next,
	}
	err := playbook.Run(config)
	return parseParams(a.Output, gv), err
}
func (a *Module) Params() interface{} {
	return map[string]interface{}{}
}

func (a *Module) Scheme() string {
	return ""
}

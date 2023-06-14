package flow

import (
	"io/ioutil"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type CronAction struct {
	Expr   string                 `yaml:"expr"`
	Tasks  []Task                 `yaml:"tasks"`
	File   string                 `yaml:"file"`
	Values map[string]interface{} `yaml:"values"`
}

func (a *CronAction) Do(conf *Config, params interface{}) (interface{}, error) {
	p := params.(CronAction)
	c := cron.New(cron.WithSeconds())
	c.AddFunc(p.Expr, func() {
		gctx := make(map[string]interface{})
		if p.File != "" {
			data, err := ioutil.ReadFile(p.File)
			if err != nil {
				logrus.Error(err)
				return
			}
			var ts []Task
			err = yaml.Unmarshal(data, &ts)
			if err != nil {
				logrus.Error(err)
				return
			}
			p.Tasks = ts
		}
		if p.Values == nil {
			p.Values = make(map[string]interface{})
		}
		gvars := NewGvars(&Vars{
			Ctx:    gctx,
			Values: p.Values,
		})
		err := NewPlaybook(p.Tasks, gvars).Run(conf)
		if err != nil {
			logrus.Error(err)
			return
		}
	})
	c.Start()
	select {}
}

func (a *CronAction) Params() interface{} {
	return CronAction{}
}

func (a *CronAction) Scheme() string {
	return ""
}

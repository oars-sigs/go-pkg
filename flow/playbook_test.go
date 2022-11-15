package flow

import (
	"testing"
)

func TestPlaybook(t *testing.T) {
	printAct := PrintAction("")
	AddCustomActions("print", &printAct)
	AddCustomActions("sum", new(Sum))
	err := Run("test.yaml")
	if err != nil {
		t.Fatal(err)
	}
}

type Sum struct {
	X int64 `yaml:"x"`
	Y int64 `yaml:"y"`
}

func (a *Sum) Do(conf *Config, params interface{}) (interface{}, error) {
	args := params.(Sum)
	return args.X + args.Y, nil
}
func (a *Sum) Params() interface{} {
	return Sum{}
}

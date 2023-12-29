package flow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Task map[string]interface{}

var inKeys = []string{
	"debug",
	"when",
	"stdout",
	"loop",
	"sleep",
	"while",
	"async",
	"await",
	"concurrency",
	"ignore_error",
	"switch",
}

func isInKey(s string) bool {
	for _, k := range inKeys {
		if k == s {
			return true
		}
	}
	return false
}

var CustomActions sync.Map

func AddCustomActions(key string, action Action) {
	CustomActions.Store(key, action)
}

type Action interface {
	Do(conf *Config, params interface{}) (interface{}, error)
	Params() interface{}
	Scheme() string
}

func (t Task) Action(conf *Config, await *gawait, vars *Gvars) (*customAction, error) {
	res := &customAction{
		gawait: await,
		vars:   vars,
		conf:   conf,
	}
	tmd, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  res,
	})
	tmd.Decode(t)
	for k, v := range t {
		if !isInKey(k) {
			if a, ok := CustomActions.Load(k); ok {
				res.a = a.(Action)
				res.params = v
			}
		}
		if k == "tasks" {
			res.params = getLoopMap(LoopRes{Item: "$.ctx.item", ItemKey: "$.ctx.itemKey"})
		}
	}
	t.SwitchTask(res)
	t.MultiTasks(res)
	t.Loop(res)
	t.While(res)
	t.DeferTask(res)
	t.Single(res)
	return res, nil
}

func (t Task) While(ctxAction *customAction) {
	s, ok := t["while"]
	if !ok {
		return
	}
	action := ctxAction.a
	m := func(conf *Config, params interface{}) (interface{}, error) {
		var res interface{}
		var err error
		i := 0
		for {
			if !when(s.(string), ctxAction.vars) {
				break
			}
			parseParams(params, ctxAction.vars.SetCtx(getLoopMap(LoopRes{Item: i})))
			res, err = ctxAction.runTask(action, conf, params)
			if err != nil {
				return nil, err
			}
			i++
		}
		return res, err
	}
	ctxAction.a = &customFuncAction{m, ctxAction.a.Params()}
}

func (t Task) Loop(ctxAction *customAction) {
	s, ok := t["loop"]
	if !ok {
		return
	}
	l := parseParams(s, ctxAction.vars)
	ls := Loop(l, ctxAction.vars)
	if len(ls) == 0 {
		return
	}

	vn := parseParams(ctxAction.Concurrency, ctxAction.vars)
	n := int64(1)
	if m, ok := vn.(int64); ok && m > 1 {
		n = m
	}
	action := ctxAction.a
	m := func(conf *Config, params interface{}) (interface{}, error) {
		var res interface{}
		var err error

		var cwg sync.WaitGroup
		cwg.Add(len(ls))
		conc := make(chan struct{}, n)
		for _, item := range ls {
			p := parseParams(params, ctxAction.vars.SetCtx(getLoopMap(item)))
			conc <- struct{}{}
			go func(p interface{}) {
				defer func() {
					cwg.Done()
					<-conc
				}()
				res, err = ctxAction.runTask(action, conf, p)
			}(p)
		}
		cwg.Wait()
		return res, err
	}
	ctxAction.a = &customFuncAction{m, ctxAction.a.Params()}
}

func (t Task) MultiTasks(ctxAction *customAction) {
	if len(ctxAction.Tasks) == 0 {
		return
	}
	m := func(conf *Config, params interface{}) (interface{}, error) {
		gv := ctxAction.vars.SetCtx(params.(map[string]interface{}))
		playbook := NewPlaybook(ctxAction.Tasks, gv)
		config := &Config{
			Workdir: ctxAction.conf.Workdir,
			Next:    playbook.Next,
		}
		err := playbook.Run(config)
		return nil, err
	}
	ctxAction.a = &customFuncAction{m, nil}
}

func (t Task) SwitchTask(ctxAction *customAction) {
	if ctxAction.Switch == nil {
		return
	}
	var res interface{}
	var err error
	m := func(conf *Config, params interface{}) (interface{}, error) {
		key := parseParams(ctxAction.Switch.Key, ctxAction.vars).(string)
		if len(ctxAction.Switch.Task) != 0 && conf.Next != nil {
			for skey, id := range ctxAction.Switch.Task {
				if skey == key {
					return conf.Next(id, ctxAction.conf, ctxAction.vars)
				}
			}
		}
		return res, err
	}
	ctxAction.a = &customFuncAction{m, nil}
}

func (t Task) Single(ctxAction *customAction) {
	if _, ok := t["loop"]; ok {
		return
	}
	if _, ok := t["while"]; ok {
		return
	}
	if _, ok := t["tasks"]; ok {
		return
	}
	if _, ok := t["switch"]; ok {
		return
	}
	action := ctxAction.a
	m := func(conf *Config, params interface{}) (interface{}, error) {
		params = parseParams(params, ctxAction.vars)
		return ctxAction.runTask(action, conf, params)
	}
	ctxAction.a = &customFuncAction{m, ctxAction.a.Params()}
}

func (t Task) DeferTask(ctxAction *customAction) {
	if len(ctxAction.DeferTasks) == 0 {
		return
	}
	m := func(conf *Config, params interface{}) (interface{}, error) {
		playbook := NewPlaybook(ctxAction.DeferTasks, ctxAction.vars)
		config := &Config{
			Workdir: ctxAction.conf.Workdir,
			Next:    playbook.Next,
		}
		err := playbook.Run(config)
		return nil, err
	}
	ctxAction.a = &customFuncAction{m, nil}
}

type customAction struct {
	a           Action
	conf        *Config
	params      interface{}
	vars        *Gvars
	gawait      *gawait
	When        string      `yaml:"when"`
	Async       string      `yaml:"async"`
	Await       []string    `yaml:"await"`
	IgnoreErr   bool        `yaml:"ignoreErr"`
	Concurrency interface{} `yaml:"concurrency"`
	Sleep       int64       `yaml:"sleep"`
	Tasks       []Task      `yaml:"tasks"`
	Switch      *Switch     `yaml:"switch"`
	Output      string      `yaml:"output"`
	DeferTasks  []Task      `yaml:"defer"`
	Debug       bool        `yaml:"debug"`
}

type Switch struct {
	Key   string            `yaml:"key"`
	Tasks map[string][]Task `yaml:"tasks"`
	Task  map[string]string `yaml:"task"`
}

func (a *customAction) Do() (interface{}, error) {
	if a.When != "" {
		if !when(a.When, a.vars) {
			return nil, nil
		}
	}
	conf := a.conf
	params := a.params
	if len(a.Await) > 0 {
		err := a.gawait.Await(a.Await, a.IgnoreErr)
		if err != nil {
			return nil, err
		}
	}
	if a.Async != "" {
		a.gawait.AddAwait(a.Async)
		go func(conf *Config, params interface{}) {
			_, err := a.runTask(a.a, conf, params)
			a.gawait.DoneAwait(a.Async, err)
		}(conf, params)
		return nil, nil
	}
	res, err := a.runTask(a.a, conf, params)
	if a.Sleep > 0 {
		time.Sleep(time.Second * time.Duration(a.Sleep))
	}
	if a.Debug {
		d, err := json.MarshalIndent(res, "", "\t")
		if err == nil {
			fmt.Println("debug:", string(d))
		} else {
			fmt.Println("debug:", res)
		}
	}
	return res, err
}

func (a *customAction) runTask(act Action, conf *Config, params interface{}) (interface{}, error) {
	aparams := act.Params()
	if _, ok := act.(*customFuncAction); !ok && aparams != nil {
		switch reflect.TypeOf(params).Kind() {
		case reflect.Map:
			md, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				TagName: "yaml",
				Result:  &aparams,
			})
			md.Decode(params)
		default:
			aparams = params
		}
	} else {
		aparams = params
	}
	res, err := act.Do(conf, aparams)
	if a.IgnoreErr {
		return res, nil
	}
	if a.Output != "" {
		a.vars.SetVar(strings.TrimPrefix(a.Output, "$."), res)
	}
	return res, err
}

type customFuncAction struct {
	m      func(conf *Config, params interface{}) (interface{}, error)
	params interface{}
}

func (a *customFuncAction) Do(conf *Config, params interface{}) (interface{}, error) {
	return a.m(conf, params)
}

func (a *customFuncAction) Params() interface{} {
	return a.params
}
func (a *customFuncAction) Scheme() string {
	return ""
}

func parseParams(ctxv interface{}, vars *Gvars) interface{} {
	if ctxv == nil {
		return ctxv
	}
	switch reflect.TypeOf(ctxv).Kind() {
	case reflect.String:
		s := ctxv.(string)
		if strings.HasPrefix(s, "$$") {
			return strings.TrimPrefix(s, "$")
		}
		if strings.HasPrefix(s, "$.") {
			p, _ := vars.GetVar(strings.TrimPrefix(s, "$."))
			return p
		}
		ss, _ := parseTpl(s, vars)
		return ss
	case reflect.Map:
		v := reflect.ValueOf(ctxv)
		res := make(map[string]interface{})
		for _, k := range v.MapKeys() {
			res[k.Interface().(string)] = parseParams(v.MapIndex(k).Interface(), vars)
		}
		return res
	case reflect.Slice:
		v := reflect.ValueOf(ctxv)
		res := make([]interface{}, 0)
		for i := 0; i < v.Len(); i++ {
			res = append(res, parseParams(v.Index(i).Interface(), vars))
		}
		return res
	case reflect.Int64:
		s := ctxv.(int64)
		return s
	case reflect.Float64:
		s := ctxv.(float64)
		return s
	}
	return ctxv
}

func parseTpl(tpl string, vars *Gvars) (string, error) {
	tmpl, err := newTpl().Parse(tpl)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	data := vars.Vars()
	vars.mutex.Lock()
	defer vars.mutex.Unlock()
	err = tmpl.Execute(&b, data)
	if err != nil {
		return "", err
	}
	v := b.String()
	if v == "<no value>" {
		return "", nil
	}
	return v, err
}

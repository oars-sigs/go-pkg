package flow

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/icza/dyno"
)

type Vars struct {
	Values map[string]interface{}
	Ctx    map[string]interface{}
}

type Gvars struct {
	data  *Vars
	mutex *sync.Mutex
}

func NewGvars(v *Vars) *Gvars {
	return &Gvars{
		data:  v,
		mutex: new(sync.Mutex),
	}
}

func (p *Gvars) SetCtx(ctx map[string]interface{}) *Gvars {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for k, v := range p.data.Ctx {
		ctx[k] = v
	}
	return &Gvars{
		data: &Vars{
			Values: p.data.Values,
			Ctx:    ctx,
		},
		mutex: p.mutex,
	}
}

func (p *Gvars) Vars() map[string]interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	res := make(map[string]interface{})
	vars := p.data
	res["values"] = mapcopy(vars.Values)
	res["ctx"] = mapcopy(vars.Ctx)
	return res
}

func (p *Gvars) GetVar(s string) (interface{}, bool) {
	vars := p.data
	p.mutex.Lock()
	defer p.mutex.Unlock()

	rs := strings.Split(s, ".")
	if len(rs) > 0 {
		var p map[string]interface{}
		switch rs[0] {
		case "values":
			p = mapcopy(vars.Values).(map[string]interface{})
		case "ctx":
			p = mapcopy(vars.Ctx).(map[string]interface{})
		default:
			return nil, false
		}
		if len(rs) == 1 {
			return p, true
		}
		for i := 1; i < len(rs); i++ {
			key := rs[i]
			if i == len(rs)-1 {
				if _, ok := p[key]; ok {
					return p[key], true
				}
			}
			if _, ok := p[key]; !ok {
				return nil, false
			}
			if nextp, ok := p[key].(map[string]interface{}); ok {
				p = nextp
				continue
			}
			if nextp, ok := p[key].(map[interface{}]interface{}); ok {
				p = mapconv(nextp)
				continue
			}
			return nil, false
		}
	}
	return nil, false
}

func (p *Gvars) SetVar(s string, value interface{}) {
	vars := p.data
	p.mutex.Lock()
	defer p.mutex.Unlock()
	rs := strings.Split(s, ".")

	if len(rs) > 1 {
		var p map[string]interface{}
		switch rs[0] {
		case "values":
			p = vars.Values
		case "ctx":
			p = vars.Ctx
		default:
			return
		}
		for i := 1; i < len(rs); i++ {
			if i == len(rs)-1 {
				p[rs[i]] = value
				return
			}
			if _, ok := p[rs[i]]; !ok {
				p[rs[i]] = make(map[string]interface{})
			}
			if _, ok := p[rs[i]].(map[string]interface{}); !ok {
				if v, ok := p[rs[i]].(map[interface{}]interface{}); ok {
					p[rs[i]] = mapconv(v)
				} else {
					p[rs[i]] = make(map[string]interface{})
				}
			}
			p = p[rs[i]].(map[string]interface{})
		}
	}
}

func mapconv(s map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range s {
		res[k.(string)] = v
	}
	return res
}

func mapcopy(v interface{}) interface{} {
	d, _ := json.Marshal(dyno.ConvertMapI2MapS(v))
	if _, ok := v.(map[interface{}]interface{}); ok {
		var res map[string]interface{}
		json.Unmarshal(d, &res)
		return dyno.ConvertMapI2MapS(res)
	}
	if _, ok := v.(map[string]interface{}); ok {
		var res map[string]interface{}
		json.Unmarshal(d, &res)
		return dyno.ConvertMapI2MapS(res)
	}
	return v
}

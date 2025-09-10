package flow

import (
	"maps"
	"strings"
	"sync"

	"github.com/icza/dyno"
)

type Vars struct {
	Values map[string]interface{}
	States map[string]interface{}
	Ctx    map[string]interface{}
	Global map[string]interface{}
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
		if _, ok := ctx[k]; !ok {
			ctx[k] = v
		}
	}
	return &Gvars{
		data: &Vars{
			Values: p.data.Values,
			States: p.data.States,
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
	res["values"] = maps.Clone(vars.Values)
	res["states"] = maps.Clone(vars.States)
	res["ctx"] = maps.Clone(vars.Ctx)
	res["global"] = maps.Clone(vars.Global)
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
		case "global":
			p = vars.Global
		case "states":
			p = vars.States
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
		case "global":
			p = vars.Global
		case "states":
			p = vars.States
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
	//d, _ := json.Marshal(dyno.ConvertMapI2MapS(v))
	return dyno.ConvertMapI2MapS(v)
	// if _, ok := v.(map[interface{}]interface{}); ok {
	// 	var res map[string]interface{}
	// 	json.Unmarshal(d, &res)
	// 	return dyno.ConvertMapI2MapS(res)
	// }
	// if _, ok := v.(map[string]interface{}); ok {
	// 	var res map[string]interface{}
	// 	json.Unmarshal(d, &res)
	// 	fmt.Println(res, dyno.ConvertMapI2MapS(v), "555555")
	// 	return dyno.ConvertMapI2MapS(res)
	// }
	// return v
}

// MergeValues Merges source and destination map, preferring values from the source map
func MergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = MergeValues(destMap, nextMap)
	}
	return dest
}

package flow

import "reflect"

func Loop(loop interface{}, vars *gvars) []LoopRes {
	switch reflect.TypeOf(loop).Kind() {
	case reflect.Slice:
		res := make([]LoopRes, 0)
		for i, item := range loop.([]interface{}) {
			res = append(res, LoopRes{
				Item:    item,
				ItemKey: i,
			})
		}
		return res
	case reflect.Map:
		res := make([]LoopRes, 0)
		v := reflect.ValueOf(loop)
		for _, k := range v.MapKeys() {
			res = append(res, LoopRes{
				Item:    v.MapIndex(k).Interface(),
				ItemKey: k,
			})
		}
		return res
	}
	return nil
}

type LoopRes struct {
	ItemKey interface{}
	Item    interface{}
}

func getLoopMap(r LoopRes) map[string]interface{} {
	return map[string]interface{}{
		"itemKey": r.ItemKey,
		"item":    r.Item,
	}
}

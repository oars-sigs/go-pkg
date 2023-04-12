package flow

func when(s string, vars *Gvars) bool {
	res, err := Eval(s, vars)
	if err != nil {
		return false
	}
	if v, ok := res.(bool); ok {
		return v
	}
	return false
}
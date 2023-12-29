package flow

type SetfaceAction string

func (a *SetfaceAction) Do(conf *Config, params interface{}) (interface{}, error) {
	return params, nil
}
func (a *SetfaceAction) Params() interface{} {
	return nil
}

func (a *SetfaceAction) Scheme() string {
	return ""
}

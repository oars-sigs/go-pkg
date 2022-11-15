package flow

type Config struct {
	Workdir string
	Next    func(id string, conf *Config, await *gawait, vars *gvars) (interface{}, error)
}

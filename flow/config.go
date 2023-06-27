package flow

type Config struct {
	Workdir string
	Next    func(id string, conf *Config, vars *Gvars) (interface{}, error)
	PTasks  []func()
}

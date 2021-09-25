package registries

var Reg *Registry

type Registry struct {
	WorkingDir string
	ExecPath   string
}

func init() {
	Reg = NewRegistry()
}

func NewRegistry() *Registry {
	return &Registry{}
}

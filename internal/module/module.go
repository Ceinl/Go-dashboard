package module

var Modules = map[string]Module{}

type Module struct{}

func init() {
	Modules["module1"] = Module{} // TODO: Change for actual module
}

package base

type Configuration struct {
	Scope			string
	Flags			[]string
	Modules			ModulesList
	Verbose			bool
}

type ModuleConfig struct{
	Path	string
	Options	map[string]string
}

type ModulesList struct {
	Processors	[]ModuleConfig
	Inspectors	[]ModuleConfig
}


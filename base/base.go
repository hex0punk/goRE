// Package base provides primitives for running gorp from the command line
package base

// Configuration holds the configuration of gorp and it is used
// when parsing the yaml config file
type Configuration struct {
	Scope			string
	Flags			[]string
	Modules			ModulesList
	Verbose			bool
}

// ModuleConfig holds the path and options for gorp modules
type ModuleConfig struct{
	Path	string
	Options	map[string]string
}

// ModulesList holds Processors and Inspectors to be used in a gorp session
type ModulesList struct {
	Processors	[]ModuleConfig
	Inspectors	[]ModuleConfig
}


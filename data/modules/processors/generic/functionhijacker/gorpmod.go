package main

import (
	"github.com/DharmaOfCode/gorp/api"
	"github.com/DharmaOfCode/gorp/modules"
	"strings"
)

type functionHijacker struct {
	Registry modules.Registry
	Options  []modules.Option
}

func (f *functionHijacker) Init() {
	f.Registry = modules.Registry{
		Name:        "functionhijacker",
		DocTypes:    []string{"Script", "XHR"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/processors/angular/functionhijacker/gorpmod.go",
		Description: "Hijacks and alters a function. The module finds the function by name.",
		Notes:       "At times a page may load scripts that have functions with the same name, in which case this would not work",
	}
	f.Options = []modules.Option{
		{
			Name:        "FunctionName",
			Value:       "",
			Required:    true,
			Description: "The name of the function to hijack",
		},
		{
			Name:        "NewBody",
			Value:       "console.log('function hijacked!')",
			Required:    true,
			Description: "The new function body",
		},
	}
}

func (f *functionHijacker) Process(webData modules.WebData) (string, error) {
	name := f.Options[0].Value
	newBody := f.Options[1].Value
	if name == "" {
		return webData.Body, nil
	}
	enableProdModeFunc, err := api.GetJsFunctionWithName(webData.Body, name)
	if err != nil || enableProdModeFunc == nil {
		// if we return an error the debugger will panic
		// and this does not warrant that
		return webData.Body, nil
	}

	return strings.Replace(webData.Body, enableProdModeFunc.Body, newBody, -1), nil
}

func (f *functionHijacker) GetRegistry() modules.Registry {
	return f.Registry
}

func (f *functionHijacker) GetOptions() []modules.Option {
	return f.Options
}

var Processor functionHijacker

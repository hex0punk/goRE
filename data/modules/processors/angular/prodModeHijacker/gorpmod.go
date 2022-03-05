package main

import (
	"github.com/hex0punk/goRE/api"
	"github.com/hex0punk/goRE/modules"
	"strings"
)

type prodModeHijacker struct {
	Registry modules.Registry
	Options  []modules.Option
}

func (p *prodModeHijacker) Init() {
	p.Registry = modules.Registry{
		Name:        "prodModeHijacker",
		DocTypes:    []string{"Script", "XHR"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/processors/angular/prodModeHijacker/gorpmod.go",
		Description: "Loads angular 2 code bundled with webpack in development mode, allowing researchers to debug dynamically from the console",
		Notes:       "",
	}
	p.Options = []modules.Option{}
}

const newProdModeFunc = `{console.log("hijacked enableProdMode function!")}`

func (p prodModeHijacker) Process(webData modules.WebData) (string, error) {
	enableProdModeFunc, err := api.GetJsFunctionWithHint(webData.Body, "Cannot enable prod mode")
	if err != nil || enableProdModeFunc == nil {
		// if we return an error the debugger will panic
		// and this does not warrant that
		return webData.Body, nil
	}
	return strings.Replace(webData.Body, enableProdModeFunc.Body, newProdModeFunc, -1), nil
}

func (p *prodModeHijacker) GetRegistry() modules.Registry {
	return p.Registry
}

func (p *prodModeHijacker) GetOptions() []modules.Option {
	return p.Options
}

var Processor prodModeHijacker

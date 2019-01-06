package main

import (
	"github.com/DharmaOfCode/gorp/modules"
	"strings"
)

type prodModeHijacker struct {
	Registry modules.Registry
	Options  []modules.Option
}

type jsFunction struct {
	Name       string
	Body       string
	Raw        string
	StartIndex int
	EndIndex   int
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

func (p *prodModeHijacker) Process(webData modules.WebData) (string, error) {
	enableProdModeFunc := getJsFunction(webData.Body, "\"Cannot enable prod mode")
	if enableProdModeFunc == nil{
		return webData.Body, nil
	}

	return strings.Replace(webData.Body, enableProdModeFunc.Body, newProdModeFunc, -1), nil
}

func getJsFunction(body string, canary string) *jsFunction {
	idx := strings.Index(body, canary)
	if idx == -1 {
		return nil
	}

	result := jsFunction{}
	// find end index for function
	// TODO: for a general function finder, this would be useless as it
	// is a naive way to parse for functions. So this needs fixed
	for i := idx; i < len(body); i++ {
		funcChar := string(body[i])
		if funcChar == "}" {
			result.EndIndex = i + 1
			break
		}
	}

	//find start index
	for i := idx; i < len(body); i-- {
		// word "function" has 8 characters
		funcWord := string(body[i-8 : i])
		if funcWord == "function" {
			result.StartIndex = i - 8
			break
		}
	}

	result.Raw = body[result.StartIndex:result.EndIndex]
	// now get the function symbol or name
	idx = strings.Index(result.Raw, "(")
	out := strings.TrimLeft(strings.TrimSuffix(result.Raw, result.Raw[idx:]), "function ")
	result.Name = strings.TrimSpace(out)

	// get the function body
	result.Body = result.Raw[strings.Index(result.Raw, "{"):]

	return &result
}

func (p *prodModeHijacker) GetRegistry() modules.Registry {
	return p.Registry
}

func (p *prodModeHijacker) GetOptions() []modules.Option {
	return p.Options
}

var Processor prodModeHijacker

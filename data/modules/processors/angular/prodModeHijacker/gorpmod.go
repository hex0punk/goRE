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
	Start      int
	End        int
	BodyStart  int
	Expression bool
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
	enableProdModeFunc := GetJsFunction(webData.Body, "\"Cannot enable prod mode")
	if enableProdModeFunc == nil {
		return webData.Body, nil
	}

	return strings.Replace(webData.Body, enableProdModeFunc.Body, newProdModeFunc, -1), nil
}

func GetJsFunction(body string, canary string) *jsFunction {
	idx := strings.Index(body, canary)
	if idx == -1 {
		return nil
	}

	result := jsFunction{}

	// find the start index for function statement/body
	// start at canary location
	// and look for a function declaration indicator
	for i := idx; i < len(body); i-- {
		if string(body[i-2:i]) == "){" {
			result.BodyStart = i - 1
			break
		} else if string(body[i-3:i]) == ") {" {
			result.BodyStart = i - 2
			break
		}
	}

	// find the start index of entire function
	// starting with word function
	for i := result.BodyStart; i < len(body); i-- {
		funcWord := string(body[i-8 : i])
		if funcWord == "function" {
			result.Start = i - 8
			break
		}
	}

	// find end index for function
	tracker := 0
	for i := result.BodyStart; i < len(body); i++ {
		if string(body[i]) == "{" {
			tracker++
		} else if string(body[i]) == "}" {
			tracker--
			if tracker == 0 {
				result.End = i + 1
				break
			}
		}
	}

	// is is declared as an expression?
	result.Expression = strings.Contains(body[result.Start:result.End], "function (") ||
		strings.Contains(body[result.Start:result.End], "function(")

	if result.Expression {
		for i := result.BodyStart; i < len(body); i-- {
			varWord := string(body[i-3 : i])
			if varWord == "var" {
				result.Start = i - 3
				break
			}
		}
	}
	result.Raw = body[result.Start:result.End]
	result.Body = body[result.BodyStart:result.End]

	// now get the function symbol or name
	var nameEnd string
	var nameBegin string
	if result.Expression {
		nameBegin = "var "
		nameEnd = "="
	} else {
		nameBegin = "function "
		nameEnd = "("
	}
	idx = strings.Index(result.Raw, nameEnd)
	out := strings.TrimLeft(strings.TrimSuffix(result.Raw, result.Raw[idx:]), nameBegin)
	result.Name = strings.TrimSpace(out)

	return &result
}

func (p *prodModeHijacker) GetRegistry() modules.Registry {
	return p.Registry
}

func (p *prodModeHijacker) GetOptions() []modules.Option {
	return p.Options
}

var Processor prodModeHijacker

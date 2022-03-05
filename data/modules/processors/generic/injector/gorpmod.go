package main

import (
	"github.com/hex0punk/goRE/api"
	"github.com/hex0punk/goRE/modules"
)

type injector struct {
	Registry modules.Registry
	Options  []modules.Option
}

func (i *injector) Init() {
	i.Registry = modules.Registry{
		Name:        "Injector",
		DocTypes:    []string{"Script", "XHR"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/processors/generic/injector/gorpmod.go",
		Description: "JS code injector.",
		Notes:       "",
	}
	i.Options = []modules.Option{
		{
			Name:        "FunctionName",
			Value:       "true",
			Required:    true,
			Description: "Function where to inject code in",
		},
		{
			Name:        "URL",
			Value:       "",
			Required:    false,
			Description: "Url of the file containing the function",
		},
		{
			Name:        "Injection",
			Value:       "",
			Required:    true,
			Description: "Code to inject",
		},
		{
			Name:        "InjectAtEnd",
			Value:       "false",
			Required:    true,
			Description: "Injects function at the end of the target function. Only useful for functions that do not have a return statement.",
		},
	}
}

func (i *injector) Process(webData modules.WebData) (string, error) {
	url := i.Options[1].Value
	if url != "" && url != webData.Url{
		return webData.Body, nil
	}

	functionName := i.Options[0].Value
	targetFunction, err := api.GetJsFunctionWithName(webData.Body, functionName)
	if err != nil || targetFunction == nil {
		return webData.Body, nil
	}

	injection := i.Options[2].Value
	atEnd := i.Options[3].Value
	var newContent string
	if atEnd ==  "true"{
		newContent = webData.Body[:targetFunction.End - 1] + injection + webData.Body[targetFunction.End -1:]
	} else {
		newContent = webData.Body[:targetFunction.BodyStart + 1] + injection + webData.Body[targetFunction.BodyStart + 1:]
	}

	return string(newContent), nil
}

func (i *injector) GetRegistry() modules.Registry {
	return i.Registry
}

func (i *injector) GetOptions() []modules.Option {
	return i.Options
}

var Processor injector

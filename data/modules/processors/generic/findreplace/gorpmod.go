package main

import (
	"github.com/DharmaOfCode/gorp/modules"
	"log"
	"strings"
)

type findreplace struct {
	Registry modules.Registry
	Options  []modules.Option
}

func (f *findreplace) Init() {
	f.Registry = modules.Registry{
		Name:        "FindReplace",
		DocTypes:    []string{"Document"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/generic/findreplace/gorpmod.go",
		Description: "Simple find replace plugin for responses.",
		Notes:       "",
	}
	f.Options = append(f.Options, modules.Option{
		Name:        "ProcessDocuments",
		Value:       "true",
		Required:    true,
		Description: "run this plugin on content of type Document",
	})

	f.Options = append(f.Options, modules.Option{
		Name:        "ProcessScripts",
		Value:       "true",
		Required:    true,
		Description: "run this plugin on content of type Script",
	})

	f.Options = append(f.Options, modules.Option{
		Name:        "Find",
		Value:       "",
		Required:    true,
		Description: "string to find",
	})

	f.Options = append(f.Options, modules.Option{
		Name:        "Replace",
		Value:       "",
		Required:    true,
		Description: "string to replace found string with",
	})

	f.Options = append(f.Options, modules.Option{
		Name:        "URL",
		Value:       "",
		Required:    false,
		Description: "URL of the file you are targeting. All files will be processed when left empty",
	})
}

func (f *findreplace) Process(webData modules.WebData) (string, error) {
	if webData.Type == "Document" && f.Options[0].Value != "true" {
		return webData.Body, nil
	}
	if webData.Type == "Script" && f.Options[1].Value != "true" {
		return webData.Body, nil
	}

	if f.Options[4].Value != "" && !strings.Contains(webData.Url, f.Options[4].Value){
		return webData.Body, nil
	}

	if !strings.Contains(webData.Body, f.Options[2].Value){
		return webData.Body, nil
	}
	log.Println("[+] findandreplace: Found something to replace!")
	return strings.Replace(webData.Body, f.Options[2].Value, f.Options[3].Value, -1), nil
}

func (f *findreplace) GetRegistry() modules.Registry {
	return f.Registry
}

func (f *findreplace) GetOptions() []modules.Option {
	return f.Options
}

var Processor findreplace

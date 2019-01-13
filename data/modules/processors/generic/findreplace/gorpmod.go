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
		DocTypes:    []string{"Document", "Script", "XHR"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/processors/generic/findreplace/gorpmod.go",
		Description: "Simple find replace plugin for responses.",
		Notes:       "",
	}
	f.Options = []modules.Option{
		{
			Name:        "ProcessDocuments",
			Value:       "true",
			Required:    true,
			Description: "run this plugin on content of type Document",
		},
		{
			Name:        "ProcessScripts",
			Value:       "true",
			Required:    true,
			Description: "run this plugin on content of type Script",
		},
		{
			Name:        "Find",
			Value:       "",
			Required:    true,
			Description: "string to find",
		},
		{
			Name:        "Replace",
			Value:       "",
			Required:    true,
			Description: "string to replace found string with",
		},
		{
			Name:        "URL",
			Value:       "",
			Required:    false,
			Description: "URL of the file you are targeting. All files will be processed when left empty",
		},
	}
}

func (f *findreplace) Process(webData modules.WebData) (string, error) {
	// This seems like a bad practice, and we could probably use modules.GetModuleOption
	// to locate each option (see apifinder module) for an example. However, I am not sure that
	// it is a better idea to iterate over a list of options every time, as gorp plugins
	// get called several times a second in some instances. Will need to determine the best approach
	// but for now it is a question of effectiveness vs. style
	if webData.Type == "Document" && f.Options[0].Value != "true" {
		return webData.Body, nil
	}
	if webData.Type == "Script" && f.Options[1].Value != "true" {
		return webData.Body, nil
	}

	url, err := modules.GetModuleOption(f.Options, "URL")
	if err != nil {
		panic(err)
	}
	if url != "" && !strings.Contains(webData.Url, url) {
		return webData.Body, nil
	}

	if !strings.Contains(webData.Body, f.Options[2].Value) {
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

package main

import (
	"github.com/DharmaOfCode/gorp/modules"
	"io/ioutil"
	"log"
	"strings"
	"github.com/DharmaOfCode/gorp/option"
)

const Delimiter = ","

type findreplace struct {
	Registry modules.Registry
	Options  []option.Option
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
	f.Options = []option.Option{
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
		{
			Name:        "NewBodyPath",
			Value:       "",
			Required:    false,
			Description: "Path for local file containing new body",
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
	urlIndex := 0
	if url.IsList(){
		urlList := url.GetAsList(Delimiter)
		if urlList != nil{
			for k, v := range urlList{
				if !strings.Contains(webData.Url, v){
					urlIndex = k
					break
				}
			}
		}

		path, err := modules.GetModuleOption(f.Options, "NewBodyPath")
		if err != nil {
			panic(err)
		}
		// Check whether we need to find and replace multiple values
		if path.IsList(){
			pathList := path.GetAsList(Delimiter)
			if pathList != nil && len(pathList) >= urlIndex {
				return f.replaceWithFile(urlList[urlIndex], pathList[urlIndex])
			}
		}
	}
	if url.Value != "" && !strings.Contains(webData.Url, url.Value) {
		return webData.Body, nil
	}
	path, err := modules.GetModuleOption(f.Options, "NewBodyPath")
	if err != nil{
		return webData.Body, nil
	}
	if path.Value != "" {
		return f.replaceWithFile(url.Value, path.Value)
	}

	if !strings.Contains(webData.Body, f.Options[2].Value) {
		return webData.Body, nil
	}

	log.Println("[+] findandreplace: Replacing content	")
	find := f.Options[2]
	replace := f.Options[3]
	replaceList := replace.GetAsList(Delimiter)
	if find.IsList() && replace.IsList(){
		for k,v := range find.GetAsList(Delimiter){
			return strings.Replace(webData.Body, v, replaceList[k], -1), nil
		}
	}

	return strings.Replace(webData.Body, f.Options[2].Value, f.Options[3].Value, -1), nil
}

func (f *findreplace) GetRegistry() modules.Registry {
	return f.Registry
}

func (f *findreplace) GetOptions() []option.Option {
	return f.Options
}

func (f *findreplace) replaceWithFile(url string, path string) (string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil{
		panic(err)}
	log.Println("[+] findandreplace: Replacing with file body")
	return string(dat), nil
}

var Processor findreplace

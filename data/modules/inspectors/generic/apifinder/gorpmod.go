package main

import (
	"github.com/DharmaOfCode/gorp/modules"
	"log"
	"os"
	"strings"
	"time"
)

type apifinder struct {
	Registry modules.Registry
	Options  []modules.Option
}

func (a *apifinder) Init() {
	a.Registry = modules.Registry{
		Name:        "APIFinder",
		DocTypes:    []string{"Document", "Script"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/generic/apifinder/gorpmod.go",
		Description: "Finds apis in javascript code and save it to a chosen file",
	}

	a.Options = []modules.Option{
		{
			Name:        "FilePath",
			Value:       "",
			Required:    true,
			Description: "The file where to save findings to",
		},
		{
			Name:        "Print",
			Value:       "true",
			Required:    true,
			Description: "When an api is found, print it to console",
		},
	}
}

func (a *apifinder) Inspect(webData modules.WebData) error {
	var f *os.File
	var err error
	////Create file if one was not provided
	fileName, err := modules.GetModuleOption(a.Options, "FilePath")
	if err != nil {
		panic(err)
	}
	if fileName == "" {
		currentTime := time.Now()
		fileName = currentTime.Format("01-02-2006") + "_apis.txt"
	}

	//We append to the existing file
	f, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	words := strings.Fields(webData.Body)
	o, err := modules.GetModuleOption(a.Options, "Print")
	if err != nil {
		panic(err)
	}
	stdOut := o == "true"
	for _, v := range words {
		if strings.Contains(v, "api/") {
			if stdOut {
				log.Println("[+] API URI:", v)
			}
			v = strings.TrimLeft(strings.TrimRight(v,`"`),`"`)
			if _, err = f.WriteString("\n[+] Possible API found in URL:" + webData.Url); err != nil {
				panic(err)
			}
			if _, err = f.WriteString("\n=========================================================="); err != nil {
				panic(err)
			}
			if _, err = f.WriteString("\n[+] API URI:" + v + "\n\n"); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

func (a *apifinder) GetRegistry() modules.Registry {
	return a.Registry
}

func (a *apifinder) GetOptions() []modules.Option {
	return a.Options
}

var Inspector apifinder

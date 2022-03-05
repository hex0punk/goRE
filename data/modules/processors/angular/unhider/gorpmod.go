package main

import (
	"fmt"
	"github.com/hex0punk/goRE/modules"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ngunhide struct {
	Registry modules.Registry
	Options  []modules.Option
}

func (n *ngunhide) Init() {
	n.Registry = modules.Registry{
		Name:        "Ng-Unhider",
		DocTypes:    []string{"Document"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/angular/unhider/gorpmod.go",
		Description: "Unhides elements hidden by angular ng-if or ngIf",
		Notes:       "This may break the functionality of some angular apps",
	}
	n.Options = []modules.Option{}
}

func (n *ngunhide) Process(webData modules.WebData) (string, error) {
	if webData.Type != "Document" || webData.Url == "http://merchant.notjet.net/header" {
		return webData.Body, nil
	}
	r := strings.NewReader(webData.Body)
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return "", err
	}

	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		//Angular 1.X
		_, ex := s.Attr("ng-if")
		if ex {
			s.SetAttr("ng-if", "true")
		}

		//Angular 2.X >
		_, ex = s.Attr("*ngIf")
		if ex {
			s.SetAttr("*ngIf", "true")
		}
	})
	res, err := doc.Html()
	if err != nil{
		fmt.Println("BAD ERROR")
	}
	fmt.Println("returning: " + res)
	
	return res, err
}

func (n *ngunhide) GetRegistry() modules.Registry {
	return n.Registry
}

func (n *ngunhide) GetOptions() []modules.Option {
	return n.Options
}

var Processor ngunhide

package main

import (
	"github.com/DharmaOfCode/gorp/modules"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ngunhide struct {
	Registry	modules.Registry
	Options		[]modules.Option
}

func (n *ngunhide) Init(){
	n.Registry = modules.Registry{
		Name: "Ng-Unhider",
		DocTypes:     []string{"Document"},
		Author: []string{"codedharma", "hex0punk"},
		Path: "./data/modules/angular/unhider/gorpmod.go",
		Description: "Unhides elements hidden by angular ng-if or ngIf",
		Notes: "This may break the functionality of some angular apps",
	}
	n.Options = []modules.Option{}
}

func (n *ngunhide) Process(body string, docType string) (string, error){
	if docType != "Document"{
		return body, nil
	}
	r := strings.NewReader(body)
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

	return doc.Html()
}

func (n *ngunhide) GetRegistry() modules.Registry{
	return n.Registry
}

func (n *ngunhide) GetOptions() []modules.Option{
	return n.Options
}

var Processor ngunhide
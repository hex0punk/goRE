package main

import (
	"github.com/DharmaOfCode/gorp/modules"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ngunhide string

var(
	registry = modules.Registry{
		Name: "Ng-Unhider",
		Author: []string{"codedharma", "hex0punk"},
		Path: "./data/modules/angular/unhider/gorpmod.go",
		Description: "Unhides elements hidden by angular ng-if or ngIf",
		Notes: "This may break the functionality of some angular apps",
	}
)

func (u ngunhide) Process(body string) (string, error){
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

func (u ngunhide) GetRegistry() modules.Registry{
	return registry
}

var Processor ngunhide
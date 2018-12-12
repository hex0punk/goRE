package main

import (
	"fmt"
	"github.com/DharmaOfCode/gorp/modules"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type unhide string

var(
	registry = modules.Registry{
		Name: "HTML-Unhider",
		Author: []string{"codedharma", "hex0punk"},
		Path: "./data/modules/generic/unhider/gorpmod.go",
		Description: "Unhides input elements from responses and adds an indicator for the name attribute",
	}
)

func (u unhide) Process(body string) (string, error){
	fmt.Println("Running unhider module...")
	r := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return "", err
	}

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		att, ex := s.Attr("type")
		//s.SetAttr("value", "TEST HERE")
		if ex && att == "hidden" {
			var v string
			v, ex = s.Attr("name")
			if ex { s.BeforeHtml("<span style='color: white; background-color:black;'>" + v +"</span>") }
			s.SetAttr("type", "")
		}
	})

	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		att, ex := s.Attr("class")
		if ex && att == "hidden" {
			s.SetAttr("class", "")
		}
	})

	return doc.Html()
}

func (u unhide) GetRegistry() modules.Registry{
	return registry
}

var Processor unhide
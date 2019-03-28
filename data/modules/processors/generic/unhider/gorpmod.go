package main

import (
	"fmt"
	"github.com/DharmaOfCode/gorp/modules"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"github.com/DharmaOfCode/gorp/option"
)

type unhide struct {
	Registry modules.Registry
	Options  []option.Option
}

func (u *unhide) Init() {
	u.Registry = modules.Registry{
		Name:        "HTML-Unhider",
		DocTypes:    []string{"Document"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/generic/unhider/gorpmod.go",
		Description: "Unhides input elements from responses and adds an indicator for the name attribute",
	}
	u.Options = []option.Option{}
}

func (u *unhide) Process(webData modules.WebData) (string, error) {
	if webData.Type != "Document" {
		return webData.Body, nil
	}
	fmt.Println("Running unhider module...")
	r := strings.NewReader(webData.Body)
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
			if ex {
				s.BeforeHtml("<span style='color: white; background-color:black;'>" + v + "</span>")
			}
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

func (u *unhide) GetRegistry() modules.Registry {
	return u.Registry
}

func (u *unhide) GetOptions() []option.Option {
	return u.Options
}

var Processor unhide

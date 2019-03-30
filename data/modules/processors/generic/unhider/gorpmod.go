package main

import (
	"fmt"
	"github.com/DharmaOfCode/gorp/modules"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type unhide struct {
	Registry modules.Registry
	Options  []modules.Option
}

func (u *unhide) Init() {
	u.Registry = modules.Registry{
		Name:        "HTML-Unhider",
		DocTypes:    []string{"Document"},
		Author:      []string{"codedharma", "hex0punk"},
		Path:        "./data/modules/generic/unhider/gorpmod.go",
		Description: "Unhides input elements from responses and adds an indicator for the name attribute",
	}
	u.Options = []modules.Option{}
}

func (u *unhide) Process(webData modules.WebData) (string, error) {
	if webData.Type != "Document" || webData.Url == "http://merchant.notjet.net/header"{
		return webData.Body, nil
	}

	r := strings.NewReader(webData.Body)
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return "", err
	}

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		att, ex := s.Attr("type")
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
	_, err = doc.Html()
	if err != nil{
		fmt.Println("BAD ERROR")
	}
	fmt.Println("returning: " + webData.Body)
	return webData.Body, err
}

func (u *unhide) GetRegistry() modules.Registry {
	return u.Registry
}

func (u *unhide) GetOptions() []modules.Option {
	return u.Options
}

var Processor unhide

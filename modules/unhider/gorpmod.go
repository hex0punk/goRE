package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type unhide string

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
			if ex { s.AfterHtml("<span style='color: white; background-color:black;'>" + v +"</span>") }
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
var Processor unhide
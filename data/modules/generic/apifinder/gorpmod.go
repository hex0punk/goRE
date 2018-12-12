package apifinder

import (
	"log"
	"strings"
)

type apifinder string

func  (a apifinder) Process(content string){
	words := strings.Fields(content)
	for _, v := range words{
		if strings.Contains(v, "/api/"){
			log.Println("[+] API URI:",  v)
		}
	}
}

var Processor apifinder
package main

import "strings"

func  findAPIs(content string){
	words := strings.Fields(content)
	for _, v := range words{
		if strings.Contains(v, "/api/"){
			log.Println("[+] API URI:",  v)
		}
	}
}

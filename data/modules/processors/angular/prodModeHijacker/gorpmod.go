package main

import (
	"fmt"
	"strings"
)

var body = `function R(t){t._nesting--,I(t)};function B(t){Xo=t};function Va(){if(Qo)throw new Error("Cannot enable prod mode after platform setup.");Yo=!1};function q(){return Qo=!0,Yo}function L(t,n,e){void 0===e&&(e=[]);var o=new Ee("Platform: "+n);`
var newProdModeFunc = `{console.log("hijacked enableProdMode function!")}`
func main(){
	idx := strings.Index(body, "\"Cannot enable prod mode")
	if idx == -1{
		return
	}

	//find end index for enableProdMode function
	funcEndIndex  := 0
	for i := idx; i < len(body); i++ {
		funcChar := string(body[i])
		if funcChar == "}"{
			funcEndIndex = i + 1
			break
		}
	}

	//find beginning index for enbleProdMod
	funcBeginIndex := 0
	for i := idx; i < len(body); i-- {
		funcWord := string(body[i-8:i])
		if funcWord == "function"{
			funcBeginIndex = i-8
			break
		}
	}

	prodModFunc := body[funcBeginIndex:funcEndIndex]
	fmt.Println(prodModFunc + "\n\n")

	// now get the function symbol or name
	idx = strings.Index(prodModFunc, "(")
	out := strings.TrimLeft(strings.TrimSuffix(prodModFunc,prodModFunc[idx:]),"function ")
	fmt.Println("out = " + out)
	funcSymbol := strings.TrimSpace(out)
	fmt.Println("func name is :" + funcSymbol)

	// get the function body
	funcBody := prodModFunc[strings.Index(prodModFunc, "{"):]
	fmt.Println("function body:" + funcBody)

	fmt.Println(strings.Replace(body, funcBody, newProdModeFunc, -1))
}
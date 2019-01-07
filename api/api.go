// Package debugger provides utilities and structs that can be used by modules
package api

import (
	"fmt"
	"strings"
)

// JsFunction holds a JS function data
type JsFunction struct {
	Name       string
	Body       string
	Raw        string
	Start      int
	End        int
	BodyStart  int
	Expression bool
}

// GetJSFunctionWithHint finds a function in a js file using a hint to locate it.
// It returns a pointer to a jsFunction object
func GetJsFunctionWithHint(body string, hint string) *JsFunction {
	// TODO: break this into smaller private functions
	// TODO: same function but find by name
	idx := strings.Index(body, hint)
	if idx == -1 {
		return nil
	}

	result := JsFunction{}

	// find the start index for function statement/body
	// start at hint location
	// and look for a function declaration indicator
	fmt.Print("finding bodystart with index: ")
	fmt.Println(idx)
	found := false
	for i := idx; i > 0; i-- {
		//fmt.Println(string(body[i-2:i]))
		if string(body[i-2:i]) == "){" {
			fmt.Println("first possible body start")
			for x := i - 1; i > 0; x--{
				if string(body[x]) == "("{
					fmt.Println("first open param")
					//check - 8 , if the word is function then steo here
					if string(body[x-8:x]) == "function" || string(body[x-9:x]) == "function "{
						fmt.Println("first found!")
						result.BodyStart = i - 1
						found = true
						break
					} else {
						//else
						//keep checking until we hit a space
						//then, if the current index - 8 = the word function then steop here
						if string(body[x]) == " " && string(body[x-8:x]) == "function"{
							result.BodyStart = i - 1
							break
						}
						continue
					}
				}
			}
			if found{
				break
			}
		} else if string(body[i-3:i]) == ") {" {
			fmt.Println("second possible body start")
			for x := i - 1; i > 0; x--{
				if string(body[x]) == "("{
					fmt.Println("Second open param")
					//check - 8 , if the word is function then steo here
					if string(body[x-8:x]) == "function" || string(body[x-9:x]) == "function "{
						fmt.Println("second found!")
						result.BodyStart = i - 2
						break
					} else {
						//else
						//keep checking until we hit a space
						//then, if the current index - 8 = the word function then steop here
						if string(body[x]) == " " && string(body[x-8:x]) == "function"{
							result.BodyStart = i - 2
							found = true
							break
						}
						continue
					}
				}
			}
			if found{
				break
			}
		}
	}

	fmt.Println(result)
	// find the start index of entire function
	// starting with word function
	for i := result.BodyStart; i < len(body); i-- {
		funcWord := string(body[i-8 : i])
		if funcWord == "function" {
			result.Start = i - 8
			break
		}
	}

	// find end index for function
	tracker := 0
	for i := result.BodyStart; i < len(body); i++ {
		if string(body[i]) == "{" {
			tracker++
		} else if string(body[i]) == "}" {
			tracker--
			if tracker == 0 {
				result.End = i + 1
				break
			}
		}
	}

	// is is declared as an expression?
	result.Expression = strings.Contains(body[result.Start:result.End], "function (") ||
		strings.Contains(body[result.Start:result.End], "function(")

	if result.Expression {
		for i := result.BodyStart; i < len(body); i-- {
			varWord := string(body[i-3 : i])
			if varWord == "var" {
				result.Start = i - 3
				break
			}
		}
	}
	result.Raw = body[result.Start:result.End]
	result.Body = body[result.BodyStart:result.End]

	// now get the function symbol or name
	var nameEnd string
	var nameBegin string
	if result.Expression {
		nameBegin = "var "
		nameEnd = "="
	} else {
		nameBegin = "function "
		nameEnd = "("
	}
	idx = strings.Index(result.Raw, nameEnd)
	out := strings.TrimLeft(strings.TrimSuffix(result.Raw, result.Raw[idx:]), nameBegin)
	result.Name = strings.TrimSpace(out)

	return &result
}
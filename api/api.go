// Package debugger provides utilities and structs that can be used by modules
package api

import "strings"

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
	for i := idx; i < len(body); i-- {
		if string(body[i-2:i]) == "){" {
			result.BodyStart = i - 1
			break
		} else if string(body[i-3:i]) == ") {" {
			result.BodyStart = i - 2
			break
		}
	}

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
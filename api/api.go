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
func GetJsFunctionWithHint(body string, hint string) (*JsFunction, error) {
	// TODO: find params as well
	// TODO: this could break if closures inside function
	idx := strings.Index(body, hint)
	if idx == -1 {
		return nil, fmt.Errorf("unable to locate function")
	}

	result := JsFunction{}

	// find the start index for function statement/body
	// start at hint location
	// and look for a function declaration indicator
	for i := idx; i > 0; i-- {
		if string(body[i-8:i]) == "function" {
			for x := i; x < len(body); x++ {
				if string(body[x]) == "{" {
					result.BodyStart = x
					break
				}
			}
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

	err := processJsFunction(&result, body)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetJsFunctionWithName(body string, name string) (*JsFunction, error) {
	signatures := []string{
		"var " + name + " = function",
		"var " + name + " =function",
		"var " + name + "= function",
		"var " + name + "=function",
		"function " + name + " (",
		"function " + name + "(",
	}
	result := JsFunction{
		Start: -1,
		Name:  name,
	}
	for _, signature := range signatures {
		result.Start = strings.Index(body, signature)
		if result.Start != -1 {
			break
		}
	}
	if result.Start == -1 {
		return nil, fmt.Errorf("unable to locate function")
	}

	//All we need to do is find the body start index
	for i := result.Start; i > 0; i++ {
		if string(body[i]) == "{" {
			result.BodyStart = i
			break
		}
	}
	err := processJsFunction(&result, body)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// This is only to avoid code duplication but it requires that
// BodyStart and Start values are found first for it to work
func processJsFunction(j *JsFunction, body string) error {
	// find end index for function
	if j.Start == 0 || j.BodyStart == 0 {
		return fmt.Errorf("JsFunction pointer must include values for Start and Body Start")
	}
	tracker := 0
	for i := j.BodyStart; i < len(body); i++ {
		if string(body[i]) == "{" {
			tracker++
		} else if string(body[i]) == "}" {
			tracker--
			if tracker == 0 {
				j.End = i + 1
				break
			}
		}
	}

	// is is declared as an expression?
	j.Expression = strings.Contains(body[j.Start:j.End], "function (") ||
		strings.Contains(body[j.Start:j.End], "function(")

	if j.Expression {
		for i := j.BodyStart; i < len(body); i-- {
			varWord := string(body[i-3 : i])
			if varWord == "var" {
				j.Start = i - 3
				break
			}
		}
	}
	j.Raw = body[j.Start:j.End]
	j.Body = body[j.BodyStart:j.End]

	// now get the function symbol or name
	if j.Name == "" {
		var nameEnd string
		var nameBegin string
		if j.Expression {
			nameBegin = "var "
			nameEnd = "="
		} else {
			nameBegin = "function "
			nameEnd = "("
		}
		idx := strings.Index(j.Raw, nameEnd)
		out := strings.TrimLeft(strings.TrimSuffix(j.Raw, j.Raw[idx:]), nameBegin)
		j.Name = strings.TrimSpace(out)
	}
	return nil
}

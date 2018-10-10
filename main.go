package main

import (
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gobs/args"
	"github.com/raff/godet"
)

type DebuggerOptions struct {
	EnableConsole 	bool
	Verbose       	bool

	AlterDocument	bool
	AlterScript		bool
}

type State struct {
	Debugger      *godet.RemoteDebugger
	Done          chan bool
	Options 	  DebuggerOptions
}


func OpenChrome(portNumber int) error {
	var chromeapp string
	chromeapp = `open -na "/Applications/Google Chrome.app" --args
			--remote-debugging-port=` + strconv.Itoa(portNumber) + `
			--window-size=1200,800
			--user-data-dir=/tmp/chrome-testing
			--auto-open-devtools-for-tabs`

	log.Println("[+] opening chrome:" + chromeapp)

	err := runCommand(chromeapp)

	return err
}

func SetupDebugger(s *State, portNumber int) {
	// Requires an opened browser running DevTools protocol
	var err error

	// Keep checking for browser with [portNumber] connection
	log.Println("Attempting to connect to browser on " + "localhost:"+strconv.Itoa(portNumber))
	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		s.Debugger, err = godet.Connect("localhost:"+strconv.Itoa(portNumber), s.Options.Verbose)
		if err == nil {
			break
		}

		log.Println("[+] Connect", err)
	}

	if err != nil {
		log.Fatal("[-] Unable to connect to the browser!")
	}
}

func SetupConsoleLogging(s *State) {
	//Log console
	log.Println("[+] Setting up console events.")
	s.Debugger.CallbackEvent("Log.entryAdded", func(params godet.Params) {
		entry := params.Map("entry")
		log.Println("LOG", entry["type"], entry["level"], entry["text"])
	})

	s.Debugger.CallbackEvent("Runtime.consoleAPICalled", func(params godet.Params) {
		l := []interface{}{"CONSOLE", params["type"].(string)}

		for _, a := range params["args"].([]interface{}) {
			arg := a.(map[string]interface{})

			if arg["value"] != nil {
				l = append(l, arg["value"])
			} else if arg["preview"] != nil {
				arg := arg["preview"].(map[string]interface{})

				v := arg["description"].(string) + "{"

				for i, p := range arg["properties"].([]interface{}) {
					if i > 0 {
						v += ", "
					}

					prop := p.(map[string]interface{})
					if prop["name"] != nil {
						v += fmt.Sprintf("%q: ", prop["name"])
					}

					v += fmt.Sprintf("%v", prop["value"])
				}

				v += "}"
				l = append(l, v)
			} else {
				l = append(l, arg["type"].(string))
			}
		}
		log.Println(l...)
	})
}

func SetupRequestInterception(s *State, requestPatterns ...godet.RequestPattern) {
	log.Println("[+] Setting up interception.")
	s.Debugger.SetRequestInterception(requestPatterns ...)
	responses := map[string]string{}
	s.Debugger.CallbackEvent("Network.requestIntercepted", func(params godet.Params) {
		iid := params.String("interceptionId")
		rtype := params.String("resourceType")
		reason := responses[rtype]
		log.Println("[+] Request intercepted for", iid, rtype, params.Map("request")["url"])
		if reason != "" {
			log.Println("  abort with reason", reason)
		}

		// Alter HTML in request response
		if s.Options.AlterDocument && rtype == "Document" && iid != "" {
			res, err := s.Debugger.SendRequest("Network.getResponseBodyForInterception", godet.Params{
				"interceptionId": iid,
			})

			if err != nil {
				log.Println("[-] Unable to get intercepted response body!")
			}

			rawAlteredResponse, err := AlterDocument(res)
			if err != nil{
				log.Println("[-] Unable to alter HTML")
			}

			if rawAlteredResponse != "" {
				log.Println("[+] Sending modified body")
				s.Debugger.ContinueInterceptedRequest(iid, godet.ErrorReason(reason), rawAlteredResponse, "", "", "", nil)
			}
		} else {
			s.Debugger.ContinueInterceptedRequest(iid, godet.ErrorReason(reason), "", "", "", "", nil)
		}
	})
}

func AlterDocument(debuggerResponse map[string]interface{}) (string, error) {
	var responseBody []byte
	if b, ok := debuggerResponse["base64Encoded"]; ok && b.(bool) {
		responseBody, _ = base64.StdEncoding.DecodeString(debuggerResponse["body"].(string))
	} else {
		responseBody = []byte(debuggerResponse["body"].(string))
	}

	alteredBody, err := processHtml(responseBody)
	if err != nil {
		return "", err
	}

	alteredHeader := "Date: " + fmt.Sprintf("%s", time.Now().Format(time.RFC3339)) + "\r\n" +
		"Connection : close\r\n" +
		"Content-Length: " + strconv.Itoa(len(alteredBody)) + "\r\n" +
		"Content-Type: text/html; charset=utf-8"

	rawAlteredResponse := base64.StdEncoding.EncodeToString([]byte("HTTP/1.1 200 OK" + "\r\n" + alteredHeader + "\r\n\r\n\r\n" + alteredBody))

	return rawAlteredResponse, nil
}

func EnableAllEvents(s *State) {
	log.Println("[+] Enabling all debugger events.")
	s.Debugger.RuntimeEvents(true)
	s.Debugger.NetworkEvents(true)
	s.Debugger.PageEvents(true) //Not used at the moment, but enabling anyways
	s.Debugger.DOMEvents(true)  //Not used at the moment, but enabling anyways
	s.Debugger.LogEvents(true)
}

func main() {
	portNumber := 9222
	s := State{}
	// This is silly, but this is just me preparing the code to use github.com/spf13/cobra
	s.Options = DebuggerOptions {
		Verbose : false,
		EnableConsole : true,
		AlterDocument: true,
		AlterScript: true,
	}

	// Launch chrome with debugging port
	err := OpenChrome(portNumber)
	if err != nil{
		log.Println("[-] Unable to start browser!")
	}

	// Get a debugger reference
	SetupDebugger(&s, portNumber)
	defer s.Debugger.Close()

	s.Done = make(chan bool)

	shouldWait := true

	//Handle connection termination and expired debugging events
	s.Debugger.CallbackEvent(godet.EventClosed, func(params godet.Params) {
		log.Println("[-] Remote Debugger connection terminated!")
		s.Done <- true
	})

	//Enable Console methods
	SetupConsoleLogging(&s)

	//Enable all debugger events
	EnableAllEvents(&s)

	//Set Network Request Interceptor Patterns for terminal logging
	htmlRequestPattern := godet.RequestPattern{
		UrlPattern:        "*",
		ResourceType:      "Document",
		InterceptionStage: "HeadersReceived",
	}

	jsRequestPattern := godet.RequestPattern{
		UrlPattern:        "*.js",
		ResourceType:      "Script",
		InterceptionStage: "HeadersReceived",
	}

	//Setup intercept event behavior
	SetupRequestInterception(&s, htmlRequestPattern, jsRequestPattern)

	//Keep this running
	if shouldWait {
		log.Println("[+] Wait for events...")
		<- s.Done
	}
}

func runCommand(commandString string) error {
	parts := args.GetArgs(commandString)
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Start()
}

func processHtml(body []byte) (string, error) {
	bodyString := string(body[:])
	r := strings.NewReader(bodyString)
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return "", err
	}
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		att, ex := s.Attr("type")
		if ex && att == "hidden" {
			s.SetAttr("type", "")
		}
	})

	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		//Angular 1.X
		_, ex := s.Attr("ng-if")
		if ex {
			s.SetAttr("ng-if", "true")
		}

		//Angular 2.X >
		_, ex = s.Attr("*ngIf")
		if ex {
			s.SetAttr("*ngIf", "true")
		}
	})

	return doc.Html()
}

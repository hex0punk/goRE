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

func main() {
	var chromeapp string
	chromeapp = `open -na "/Applications/Google Chrome.app" --args
			--remote-debugging-port=9222
			--window-size=1200,800
			--user-data-dir=/tmp/chrome-testing
			--auto-open-devtools-for-tabs`

	log.Println("[+] opening chrome:" + chromeapp)
	if err := runCommand(chromeapp); err != nil {
		log.Println("[-] Unable to start browser!")
	}

	//Connect to chrome
	var remote *godet.RemoteDebugger
	var err error

	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		remote, err = godet.Connect("localhost:9222", false)
		if err == nil {
			break
		}

		log.Println("[+] Connect", err)
	}

	if err != nil {
		log.Fatal("[-] Unable to connect to the browser!")
	}

	defer remote.Close()

	done := make(chan bool)
	shouldWait := true

	var pwait chan bool

	//Handle connection termination and expired debugging events
	remote.CallbackEvent(godet.EventClosed, func(params godet.Params) {
		log.Println("[-] RemoteDebugger connection terminated!")
		done <- true
	})

	remote.CallbackEvent("Emulation.virtualTimeBudgetExpired", func(params godet.Params) {
		pwait <- true
	})

	//Log console
	remote.CallbackEvent("Log.entryAdded", func(params godet.Params) {
		entry := params.Map("entry")
		log.Println("LOG", entry["type"], entry["level"], entry["text"])
	})

	remote.CallbackEvent("Runtime.consoleAPICalled", func(params godet.Params) {
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

	remote.RuntimeEvents(true)
	remote.NetworkEvents(true)
	remote.PageEvents(true) //Not used at the moment, but enabling anyways
	remote.DOMEvents(true) //Not used at the moment, but enabling anyways
	remote.LogEvents(true) //Not used at the moment, but enabling anyways

	//Set Network Request Interception for terminal logging
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
	remote.SetRequestInterception(htmlRequestPattern, jsRequestPattern)
	responses := map[string]string{}
	remote.CallbackEvent("Network.requestIntercepted", func(params godet.Params) {
		iid := params.String("interceptionId")
		rtype := params.String("resourceType")
		reason := responses[rtype]
		log.Println("[+] Request intercepted for", iid, rtype, params.Map("request")["url"])
		if reason != "" {
			log.Println("  abort with reason", reason)
		}

		var newBody string
		var newHeaders string

		//Test intercepting HTML
		if rtype == "Document" && iid != "" {
			log.Println("[+] Request intercepted for", iid, rtype, params.Map("request")["url"])

			//Get response body for interception
			res, err := remote.SendRequest("Network.getResponseBodyForInterception", godet.Params{
				"interceptionId": iid,
			})

			var responseBody []byte
			if err != nil {
				log.Println("[-] Unable to get intercepted response body!")
			} else if b, ok := res["base64Encoded"]; ok && b.(bool) {
				responseBody, _ = base64.StdEncoding.DecodeString(res["body"].(string))
			} else {
				responseBody = []byte(res["body"].(string))
			}

			newBody, err = processHtml(responseBody)
			if err !=  nil {
				log.Println("[-] Unable to process HTML")
			}

			newHeaders = "Date: " + fmt.Sprintf("%s", time.Now().Format(time.RFC3339)) + "\r\n" +
				"Connection : close\r\n" +
				"Content-Length: " + strconv.Itoa(len(newBody)) + "\r\n" +
				"Content-Type: text/html; charset=utf-8"

		}

		if newBody != "" && newHeaders != ""{
			log.Println("[+] Sending modified body")
			rawResponse := base64.StdEncoding.EncodeToString([]byte("HTTP/1.1 200 OK" + "\r\n" + newHeaders + "\r\n\r\n\r\n" + newBody))
			remote.ContinueInterceptedRequest(iid, godet.ErrorReason(reason), rawResponse, "", "", "", nil)
		} else {
			remote.ContinueInterceptedRequest(iid, godet.ErrorReason(reason), "", "", "", "", nil)
		}
	})

	//Keep this running
	if shouldWait {
		log.Println("[+] Wait for events...")
		<-done
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

	if err != nil{
		return "", err
	}
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		att, ex := s.Attr("type")
		if ex && att == "hidden"{
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


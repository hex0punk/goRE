package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/wirepair/gcd"
	"github.com/wirepair/gcd/gcdapi"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	testListener   net.Listener
	testPath       string
	testDir        string
	testPort       string
	testServerAddr string
)

type DebuggerOptions struct {
	EnableConsole bool
	Verbose       bool

	AlterDocument bool
	AlterScript   bool
}

type State struct {
	Debugger *gcd.Gcd
	Done     chan bool
	Options  DebuggerOptions
	Target 	 *gcd.ChromeTarget
}

var testStartupFlags = []string{"--disable-new-tab-first-run", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"}

func init() {
	flag.StringVar(&testPath, "chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "path to chrome")
	flag.StringVar(&testDir, "dir", "/tmp/chrome-testing", "user directory")
	flag.StringVar(&testPort, "port", "9222", "Debugger port")
}

func AlterDocument(debuggerResponse string, headers map[string]interface{}) (string, error) {
	alteredBody, err := processHtml(debuggerResponse)
	if err != nil {
		return "", err
	}

	//gzip := false
	alteredHeader := ""
	for k, v := range headers {
		switch strings.ToLower(k) {
		case "content-length":
			v = strconv.Itoa(len(alteredBody))
			break
		case "date":
			v = fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
			break
			//case "content-encoding":
			//	ce := v.(string)
			//	gzip = ce == "gzip"
			//	break
		}
		alteredHeader += k + ": " + v.(string) + "\r\n"
	}
	alteredHeader += "\r\n"

	// This does not seem needed at the moment
	//if gzip {
	//	alteredBody = gZipCompress(alteredBody)
	//}

	rawAlteredResponse := base64.StdEncoding.EncodeToString([]byte("HTTP/1.1 200 OK"+"\r\n"+alteredHeader+alteredBody))

	return rawAlteredResponse, nil
}

//Enable request interception using the specific requestPatterns
func SetupRequestInterception(s *State, requestPatterns []*gcdapi.NetworkRequestPattern) {
	s.Target.Network.SetRequestInterception(requestPatterns)
	s.Target.Subscribe("Network.requestIntercepted", func(target *gcd.ChromeTarget, v []byte) {

		msg := &gcdapi.NetworkRequestInterceptedEvent{}
		err := json.Unmarshal(v, msg)
		if err != nil {
			log.Fatalf("error unmarshalling event data: %v\n", err)
		}
		log.Println("Method: %s\n", msg.Method)
		iid := msg.Params.InterceptionId
		rtype := msg.Params.ResourceType
		reason := msg.Params.ResponseErrorReason
		url := msg.Params.Request.Url
		responseHeaders := msg.Params.ResponseHeaders

		log.Println("[+] Request intercepted for", iid, rtype, url)
		if reason != "" {
			log.Println("[-] Abort with reason", reason)
		}

		if s.Options.AlterDocument && rtype == "Document" && iid != "" {
			res, encoded, err := target.Network.GetResponseBodyForInterception(iid)
			if err != nil {
				log.Println("[-] Unable to get intercepted response body!")
			}
			if encoded{
				res, err = decodeBase64Response(res)
				if err != nil {
					log.Println("[-] Unable to decode body!")
				}
			}

			rawAlteredResponse, err := AlterDocument(res, responseHeaders)
			if err != nil {
				log.Println("[-] Unable to alter HTML")
			}

			if rawAlteredResponse != "" {
				log.Println("[+] Sending modified body")

				_, err := target.Network.ContinueInterceptedRequest(iid, reason, rawAlteredResponse, "", "", "", nil, nil)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			target.Network.ContinueInterceptedRequest(iid, reason, "", "", "", "", nil, nil)
		}
	})
}

func main(){
	s := State{}
	// This is silly, but this is just me preparing the code to use github.com/spf13/cobra
	s.Options = DebuggerOptions{
		Verbose:       false,
		EnableConsole: true,
		AlterDocument: true,
		AlterScript:   true,
	}

	s.Debugger = startGcd()
	defer s.Debugger.ExitProcess()

	s.Target = startTarget(s.Debugger)
	//Create a channel to be able to signal a termination to our Chrome connection
	s.Done = make(chan bool)
	shouldWait := true


	//Set Network Request Interceptor Patterns for terminal logging
	htmlRequestPattern := gcdapi.NetworkRequestPattern {
		UrlPattern:        "*",
		ResourceType:      "Document",
		InterceptionStage: "HeadersReceived",
	}

	jsRequestPattern := gcdapi.NetworkRequestPattern{
		UrlPattern:        "*.js",
		ResourceType:      "Script",
		InterceptionStage: "HeadersReceived",
	}

	reqPattern := []*gcdapi.NetworkRequestPattern{&jsRequestPattern, &htmlRequestPattern}

	SetupRequestInterception(&s, reqPattern)

	if shouldWait {
		log.Println("[+] Wait for events...")
		<-s.Done
	}
}

func startGcd() *gcd.Gcd {
	testDir = "/tmp/chrome-testing"
	testPort = "9222"
	debugger := gcd.NewChromeDebugger()
	debugger.AddFlags(testStartupFlags)
	debugger.StartProcess(testPath, testDir, testPort)
	return debugger
}

func startTarget(debugger *gcd.Gcd) *gcd.ChromeTarget {
	target, err := debugger.NewTab()
	if err != nil {
		log.Fatalf("error getting new tab: %s\n", err)
	}

	// TODO: set based on verbose flag
	target.DebugEvents(false)
	target.DOM.Enable()
	target.Console.Enable()
	target.Page.Enable()
	target.Debugger.Enable()
	//This does not seem right, but will leave this here for now
	target.Network.Enable(999999,999999,999999)
	return target

}

func decodeBase64Response(res string) (string, error) {
	l, err := base64.StdEncoding.DecodeString(res)
	if err != nil{
		return "", err
	}

	return string(l[:]), nil
}


func processHtml(body string) (string, error) {
	r := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return "", err
	}
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		att, ex := s.Attr("type")
		s.SetAttr("value", "TEST HERE")
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

func gZipCompress(content string) string {
	var b bytes.Buffer
	//btw 4 and 5
	gz, err := gzip.NewWriterLevel(&b, -1)
	if err != nil {
		panic(err)
	}
	if _, err := gz.Write([]byte(content)); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	return b.String()
}
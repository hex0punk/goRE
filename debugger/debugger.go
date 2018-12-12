package debugger

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/DharmaOfCode/gorp/modules"
	"github.com/wirepair/gcd"
	"github.com/wirepair/gcd/gcdapi"
	"log"
	"strconv"
	"strings"
	"time"
)

type Debugger struct {
	ChromeProxy *gcd.Gcd
	Done     		chan bool
	Options  		Options
	Target 	 		*gcd.ChromeTarget
	Modules 		modules.Modules
}

type Options struct {
	EnableConsole bool
	Verbose       bool
	AlterDocument bool
	AlterScript   bool
	Scope 		  string
}

func (d *Debugger) StartTarget() {
	target, err := d.ChromeProxy.NewTab()
	if err != nil {
		log.Fatalf("error getting new tab: %s\n", err)
	}

	// TODO: set based on verbose flag
	target.DebugEvents(d.Options.Verbose)
	target.DOM.Enable()
	target.Console.Enable()
	target.Page.Enable()
	target.Debugger.Enable()
	networkParams := &gcdapi.NetworkEnableParams{
		MaxTotalBufferSize:    -1,
		MaxResourceBufferSize: -1,
	}
	if _, err := target.Network.EnableWithParams(networkParams); err != nil {
		log.Fatal("[-] Error enabling network!")
	}
	d.Target = target
}

// Enable request interception using the specific requestPatterns
func (d *Debugger) SetupRequestInterception(params *gcdapi.NetworkSetRequestInterceptionParams) {
	log.Println("[+] Setting up request interception")
	if _, err := d.Target.Network.SetRequestInterceptionWithParams(params); err != nil {
		log.Println("[-] Unable to setup request interception!", err)
	}

	d.Target.Subscribe("Network.requestIntercepted", func(target *gcd.ChromeTarget, v []byte) {

		msg := &gcdapi.NetworkRequestInterceptedEvent{}
		err := json.Unmarshal(v, msg)
		if err != nil {
			log.Fatalf("error unmarshalling event data: %v\n", err)
		}
		iid := msg.Params.InterceptionId
		rtype := msg.Params.ResourceType
		reason := msg.Params.ResponseErrorReason
		responseHeaders := msg.Params.ResponseHeaders

		if msg.Params.IsNavigationRequest{
			log.Print("\n\n\n\n")
			log.Println("[?] Navigation REQUEST")
		}
		log.Println("[+] Request intercepted for", msg.Params)
		if reason != "" {
			log.Println("[-] Abort with reason", reason)
		}

		if rtype == "Script" && iid != ""{
			res, encoded, err := d.Target.Network.GetResponseBodyForInterception(iid)
			if err != nil {
				log.Println("[-] Unable to get intercepted response body!", err.Error())
				target.Network.ContinueInterceptedRequest(iid, reason, "", "", "", "", nil, nil)
			} else {
				if encoded{
					res, err = decodeBase64Response(res)
					if err != nil {
						log.Println("[-] Unable to decode body!")
					}
				}
				go findAPIs(res)
			}
		}

		if d.Options.AlterDocument && rtype == "Document" && iid != "" {
			res, encoded, err := d.Target.Network.GetResponseBodyForInterception(iid)
			if err != nil {
				log.Println("[-] Unable to get intercepted response body!", err.Error())
				target.Network.ContinueInterceptedRequest(iid, reason, "", "", "", "", nil, nil)
			} else {
				if encoded{
					res, err = decodeBase64Response(res)
					if err != nil {
						log.Println("[-] Unable to decode body!")
					}
				}

				rawAlteredResponse, err := d.AlterDocument(res, responseHeaders)
				if err != nil {
					log.Println("[-] Unable to alter HTML")
				}

				if rawAlteredResponse != "" {
					log.Print("[+] Sending modified body\n\n\n")

					_, err := d.Target.Network.ContinueInterceptedRequest(iid, reason, rawAlteredResponse, "", "", "", nil, nil)
					if err != nil {
						log.Println(err)
					}
				}
			}
		} else {
			d.Target.Network.ContinueInterceptedRequest(iid, reason, "", "", "", "", nil, nil)
		}
	})
}

func (d *Debugger) AlterDocument(debuggerResponse string, headers map[string]interface{}) (string, error) {
	alteredBody, err := d.processHtml(debuggerResponse)
	if err != nil {
		return "", err
	}

	alteredHeader := ""
	for k, v := range headers {
		switch strings.ToLower(k) {
		case "content-length":
			v = strconv.Itoa(len(alteredBody))
			break
		case "date":
			v = fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
			break
		}
		alteredHeader += k + ": " + v.(string) + "\r\n"
	}
	alteredHeader += "\r\n"

	rawAlteredResponse := base64.StdEncoding.EncodeToString([]byte("HTTP/1.1 200 OK"+"\r\n"+alteredHeader+alteredBody))

	return rawAlteredResponse, nil
}

func decodeBase64Response(res string) (string, error) {
	l, err := base64.StdEncoding.DecodeString(res)
	if err != nil{
		return "", err
	}

	return string(l[:]), nil
}

func (d *Debugger) processHtml(body string) (string, error) {
	result := body
	var err error
	for _, v := range d.Modules.Processors{
		result, err = v.Process(result)
		if err != nil {
			return "", err
		}
	}
	return result, nil
}

// This needs to be actored out of here
func findAPIs(content string){
	words := strings.Fields(content)
	for _, v := range words{
		if strings.Contains(v, "/api/"){
			log.Println("[+] API URI:",  v)
		}
	}
}



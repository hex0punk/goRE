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
	Scope 		  string
}

func (d *Debugger) StartTarget() {
	target, err := d.ChromeProxy.NewTab()
	if err != nil {
		log.Fatalf("error getting new tab: %s\n", err)
	}

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
		reason := msg.Params.ResponseErrorReason
		rtype := msg.Params.ResourceType
		responseHeaders := msg.Params.ResponseHeaders
		url := msg.Params.Request.Url

		if msg.Params.IsNavigationRequest{
			log.Print("\n\n\n\n")
			log.Println("[?] Navigation REQUEST")
		}
		log.Println("[+] Request intercepted for", url)
		if reason != "" {
			log.Println("[-] Abort with reason", reason)
		}

		if iid != "" {
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
				webData := modules.WebData{
					Body: res,
					Headers: responseHeaders,
					Type: rtype,
					Url:url,
				}
				d.InspectDocument(webData)

				if rtype != ""{
					rawAlteredResponse, err := d.AlterDocument(webData)
					if err != nil {
						log.Println("[-] Unable to alter HTML")
					}

						log.Print("[+] Sending modified body\n\n\n")

						_, err = d.Target.Network.ContinueInterceptedRequest(iid, reason, rawAlteredResponse, "", "", "", nil, nil)
						if err != nil {
							log.Println(err)
						}
				} else {
					d.Target.Network.ContinueInterceptedRequest(iid, reason, "", "", "", "", nil, nil)
				}
			}
		} else {
			d.Target.Network.ContinueInterceptedRequest(iid, reason, "", "", "", "", nil, nil)
		}
	})
}

func (d *Debugger) AlterDocument(data modules.WebData) (string, error) {
	alteredBody, err := d.processBody(data)
	if err != nil {
		return "", err
	}

	alteredHeader := ""
	for k, v := range data.Headers {
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

func (d *Debugger) InspectDocument(webData modules.WebData){
	//TODO: abstract this as a debugger function
	for _, v := range d.Modules.Inspectors{
		//TODO call all inspectors as goroutines
		err := v.Inspect(webData)
		if err != nil {
			log.Println("[+] Inspector error: " + v.Registry.Name)
		}
	}
}

func decodeBase64Response(res string) (string, error) {
	l, err := base64.StdEncoding.DecodeString(res)
	if err != nil{
		return "", err
	}

	return string(l[:]), nil
}

func (d *Debugger) processBody(data modules.WebData) (string, error) {
	result := data.Body
	var err error
	for _, v := range d.Modules.Processors{
		result, err = v.Process(data)
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



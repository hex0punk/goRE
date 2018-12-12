package main

import (
	"flag"
	"fmt"
	"github.com/DharmaOfCode/gorp/debugger"
	"github.com/DharmaOfCode/gorp/modules"
	"github.com/wirepair/gcd"
	"github.com/wirepair/gcd/gcdapi"
	"log"
	"os"
)

type State struct {
	Debugger	debugger.Debugger
	Modules     modules.Modules
}

var (
	testPath       string
	testDir        string
	testPort       string
)

var testStartupFlags = []string{"-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"}
//var inspectors = []string{"./modules/generic/apifinder"}
var processors = []string{"./data/modules/generic/unhider/", "./data/modules/angular/unhider/"}
var scope = "zomato.com"

func init() {
	flag.StringVar(&testPath, "chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "path to chrome")
	flag.StringVar(&testDir, "dir", "/tmp/chrome-testing", "user directory")
	flag.StringVar(&testPort, "port", "9222", "Debugger port")
}

func main(){
	var err error
	s := State{}

	// Load the modules
	s.Modules = modules.Modules{}
	err = s.Modules.InitProcessors(processors)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	// Setup the debugger
	s.Debugger = debugger.Debugger{
		Modules: s.Modules,
	}
	s.Debugger.Options = debugger.Options{
		Verbose:       false,
		EnableConsole: true,
		AlterDocument: true,
		AlterScript:   true,
	}

	// TODO: This should be abstracted in the debugger struct
	s.Debugger.ChromeProxy = startGcd()
	defer s.Debugger.ChromeProxy.ExitProcess()

	s.Debugger.StartTarget()
	//Create a channel to be able to signal a termination to our Chrome connection
	s.Debugger.Done = make(chan bool)
	shouldWait := true

	patterns := make([]*gcdapi.NetworkRequestPattern, 2)
	patterns[0] = &gcdapi.NetworkRequestPattern{
		UrlPattern: "*" + s.Debugger.Options.Scope + "/*",
		ResourceType: "Document",
		InterceptionStage: "HeadersReceived",
	}
	patterns[1] = &gcdapi.NetworkRequestPattern{
		UrlPattern:        "*" + s.Debugger.Options.Scope + "*.js",
		ResourceType:      "Script",
		InterceptionStage: "HeadersReceived",
	}
	interceptParams := &gcdapi.NetworkSetRequestInterceptionParams{Patterns: patterns}

	s.Debugger.SetupRequestInterception(interceptParams)

	if shouldWait {
		log.Println("[+] Wait for events...")
		<-s.Debugger.Done
	}
}

// TODO: Move this to debugger
func startGcd() *gcd.Gcd {
	testDir = "/tmp/chrome-testing"
	testPort = "9222"
	debugger := gcd.NewChromeDebugger()
	//debugger.DeleteProfileOnExit()
	debugger.AddFlags(testStartupFlags)
	debugger.StartProcess(testPath, testDir, testPort)
	return debugger
}
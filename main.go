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
	"plugin"
)

type State struct {
	Debugger	debugger.Debugger
}

var (
	testPath       string
	testDir        string
	testPort       string
)

var testStartupFlags = []string{"-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"}
//var premodules = []string{"./modules/generic/apifinder"}
var processormodules = []string{"./modules/unhider/", "./modules/angular/unhider/"}
var scope = "zomato.com"

func init() {
	flag.StringVar(&testPath, "chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "path to chrome")
	flag.StringVar(&testDir, "dir", "/tmp/chrome-testing", "user directory")
	flag.StringVar(&testPort, "port", "9222", "Debugger port")
}

//TODO: this should take a refrence for the container where to put the loaded
// modules and a list of path as string array. That way we can call this function
// for preprocessors and  postprocessors
func LoadProcessors(s *State, procList []string) error{
	// TODO: Put this in a loop
	//Debugger.Processors = make([]modules.Processor, 2)
	for _, v := range procList{
		fmt.Println("[+] Loading module: " + v)
		mod := v + "gorpmod.so"
		plug, err := plugin.Open(mod)
		if err != nil {
			return err
		}

		// look up a symbol (an exported function or variable)
		// in this case, variable Greeter
		symProcessor, err := plug.Lookup("Processor")
		if err != nil {
			return err
		}

		// Assert that loaded symbol is of a desired type
		// in this case interface type Greeter (defined above)
		var processor modules.Processor
		//processor = new(modules.Processor)
		processor, ok := symProcessor.(modules.Processor)
		if !ok {
			fmt.Println("unexpected type from module symbol")
			return err
		}
		s.Debugger.Processors = append(s.Debugger.Processors, processor)
	}

	return nil
}

func main(){
	var err error
	s := State{}

	s.Debugger = debugger.Debugger{}
	err = LoadProcessors(&s, processormodules)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	s.Debugger.Options = debugger.Options{
		Verbose:       false,
		EnableConsole: true,
		AlterDocument: true,
		AlterScript:   true,
	}

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
package main

import (
	"flag"
	"fmt"
	"github.com/DharmaOfCode/gorp/base"
	"github.com/DharmaOfCode/gorp/debugger"
	"github.com/DharmaOfCode/gorp/modules"
	"github.com/DharmaOfCode/gorp/server"
	"github.com/spf13/viper"
	"github.com/wirepair/gcd"
	"github.com/wirepair/gcd/gcdapi"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// State identifies the state of a gorp session.
type State struct {
	Debugger         debugger.Debugger // Debugger object
	Modules          modules.Modules   //Selected modules
	ModPath          string            // Module path
	Run              bool              // Whether to run a session
	GetInfo          bool              // Get module information
	RecompileModules bool
	EnableUI		 bool
}

var (
	cfgFile string
	config  *base.Configuration

	chromePath string
	dumpDir    string
	debugPort  string
)

const modulesFolder = "./data/modules/"

// Init Initializes required settings
func init() {
	flag.StringVar(&chromePath, "chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "path to chrome")
	flag.StringVar(&dumpDir, "dir", "/tmp/chrome-testing", "user directory")
	flag.StringVar(&debugPort, "port", "9222", "Debugger port")
}

// ParseCmdLine parses the command line options
func ParseCmdLine() *State {
	s := State{}
	// TODO: Config should be part of state
	flag.StringVar(&cfgFile, "c", "", "configuration file path")
	flag.StringVar(&s.ModPath, "m", "", "path of module to get info for")
	flag.BoolVar(&s.Run, "r", true, "run gorp")
	flag.BoolVar(&s.GetInfo, "i", false, "run gorp")
	flag.BoolVar(&s.RecompileModules, "p", false, "recompile all plugins")

	flag.BoolVar(&s.EnableUI, "u", false, "Run wasm based UI (Experimental")

	flag.Parse()
	return &s
}

// RecompileModules recompiles all gorp plugins
func RecompileModules() {
	folders, err := ioutil.ReadDir(modulesFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range folders {
		pluginPath := modulesFolder + f.Name()
		err = filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
			if containsGorpPlugin(path) {
				log.Println("[+] recompiling plugin module: " + filepath.Base(path))
				out, err := exec.Command("go", "build", "-buildmode=plugin", "-o", path+"/gorpmod.so", path+"/gorpmod.go").
					Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%s\n", out)
			}
			return nil
		})
	}
}

// RunGorp runs gorp
func RunGorp(s *State) {
	initConfig()
	var err error

	// Load the modules
	s.Modules = modules.Modules{}
	err = s.Modules.InitProcessors(config.Modules.Processors)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.Modules.InitInspectors(config.Modules.Inspectors)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Setup the debugger
	s.Debugger = debugger.Debugger{
		Modules: s.Modules,
	}

	s.Debugger.Options = debugger.Options{
		Verbose:       config.Verbose,
		EnableConsole: true,
		LogFile:  "./logs/testlogs.txt",
	}
	s.Debugger.SetupFileLogger()
	s.Debugger.XHRBreakPoints = config.XHRBreakPoints

	// TODO: This should be abstracted out in the debugger struct
	s.Debugger.ChromeProxy = startGcd()
	defer s.Debugger.ChromeProxy.ExitProcess()

	s.Debugger.StartTarget()
	//Create a channel to be able to signal a termination to our Chrome connection
	s.Debugger.Done = make(chan bool)
	shouldWait := true

	//Default is everything!
	docPattern := "*"
	jsPattern := "*"
	xhrPattern := "*"
	if config.Scope != "" {
		docPattern = "*" + config.Scope + "/*"
		jsPattern = "*" + config.Scope + "*.js"
		xhrPattern = "*" + config.Scope + "/*"
	}
	patterns := []*gcdapi.NetworkRequestPattern{
		{
			UrlPattern:        docPattern,
			ResourceType:      "Document",
			InterceptionStage: "HeadersReceived",
		},
		{
			UrlPattern:        jsPattern,
			ResourceType:      "Script",
			InterceptionStage: "HeadersReceived",
		},
		{
			UrlPattern:        xhrPattern,
			ResourceType:      "XHR",
			InterceptionStage: "HeadersReceived",
		},
		{
			UrlPattern:        "*" + config.Scope + "*.swf",
			ResourceType:      "Other",
			InterceptionStage: "HeadersReceived",
		},
	}

	interceptParams := &gcdapi.NetworkSetRequestInterceptionParams{Patterns: patterns}

	s.Debugger.SetupRequestInterception(interceptParams)
	s.Debugger.SetupDOMDebugger()
	//Now setup script injector
	if config.Script != nil{
		//scripts, err := GetUserScripts()
		if err != nil{
			log.Println("[-] Error setting up script injector")
		}

		s.Debugger.UpdateScriptsOnLoad(config.Script.Path)
		//if config.Script.Source == ""{
		//	s.Debugger.InjectScriptAsPageObject(&scripts)
		//} else {
		//	s.Debugger.InjectScriptAsRuntime(&scripts, &config.Script.Source)
		//}
	}
	if shouldWait {
		log.Println("[+] Waiting for events...")

		<-s.Debugger.Done
	}
}

func GetUserScripts() (string, error) {
	s, err := ioutil.ReadFile(config.Script.Path) // just pass the file name
	if err != nil {
		return "", err
	}
	// Append init function.
	// TODO: Yes, I should make it a const, at least
	scripts := "setTimeout(function() { gorp(); }, 2000);\n" + string(s)
	return scripts , nil
}

// Gets and prints the information for any given module
func GetModInfo(s *State) {
	s.Modules = modules.Modules{}
	if strings.Contains(s.ModPath, "processors") {
		p, err := s.Modules.GetProcessor(s.ModPath)
		if err != nil {
			log.Println("[+] Unable to find processor " + s.ModPath)
		} else {
			p.ShowInfo()
		}
	} else if strings.Contains(s.ModPath, "inspectors") {
		i, err := s.Modules.GetInspector(s.ModPath)
		if err != nil {
			log.Println("[+] Unable to find processor " + s.ModPath)
		} else {
			i.ShowInfo()
		}
	} else {
		log.Println("[+] Unable to find module " + s.ModPath)
	}

	fmt.Println(s.ModPath)
}


func main() {
	s := ParseCmdLine()
	if s.EnableUI{
		log.Println("[+] Setting up wasm server")
		go server.Serve()
	}
	if s.GetInfo {
		GetModInfo(s)
	} else if s.RecompileModules {
		RecompileModules()
	} else {
		RunGorp(s)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find in base
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&config)
}

// TODO: Move this to debugger
func startGcd() *gcd.Gcd {
	dumpDir = "/tmp/chrome-testing"
	debugPort = "9222"
	chromeDebugger := gcd.NewChromeDebugger()
	//debugger.DeleteProfileOnExit()
	chromeDebugger.AddFlags(config.Flags)
	err := chromeDebugger.StartProcess(chromePath, dumpDir, debugPort)
	if err != nil {
		panic(fmt.Errorf("Fatal error starting chrome debugger: %s \n", err))
		os.Exit(1)
	}

	return chromeDebugger
}

func containsGorpPlugin(path string) bool {
	if _, err := os.Stat(path + "/gorpmod.go"); err == nil {
		return true
	}
	return false
}

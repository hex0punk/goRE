package cli

import (
	"fmt"
	"github.com/DharmaOfCode/gorp/pkg/core"
	"github.com/DharmaOfCode/gorp/pkg/debugger"
	"github.com/DharmaOfCode/gorp/pkg/processors"
	"github.com/abiosoft/readline"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var serverLog *os.File
var cdb debugger.Debugger
var prompt *readline.Instance
var shellCompleter *readline.PrefixCompleter
var shellMenuContext = "main"


func Shell(){
	p, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31mGorp»\033[0m ",
		HistoryFile:         "/tmp/readline.tmp",
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
	})

	if err != nil {
		color.Red("[!]There was an error with the provided input")
		color.Red(err.Error())
	}

	prompt = p
	defer prompt.Close()

	log.SetOutput(prompt.Stderr())

	for {
		line, err := prompt.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		cmd := strings.Fields(line)

		if len(cmd) > 0 {
			switch shellMenuContext {
			case "main":
				switch cmd[0] {
				case "?", "help":
					menuHelpMain()
				case "quit", "exit":
					exit()
				}
			case "add":
				menuAdd(cmd[1:])
			case "configure":
				menuConfigure(cmd[1:])
			case "":
			default:
				message("info", "Executing system command...")
				if len(cmd) > 1 {
					executeCommand(cmd[0], cmd[1:])
				} else {
					var x []string
					executeCommand(cmd[0], x)
				}
			}
		}
	}
}


func menuAdd(cmd []string) {
	if len(cmd) > 0 {
		switch cmd[0] {
		case "processor":
			if len(cmd) > 1 {
				addProcessor(cmd[1])
			} else {
				message("warning", "invalid processor")
			}
		case "":
		default:
			color.Yellow("[-]Invalid 'add' command") //TODO: Add function for warnings
		}

	}
}

func menuConfigure(cmd []string){
	if len(cmd) > 0 {
		switch cmd[0] {
		case "processor":
			if len(cmd) > 1 {
				menuConfigureProcessor(cmd[1])
			} else {
				message("warning", "invalid processor")
			}
		case "":
		default:
			color.Yellow("[-]Invalid 'configure' command") //TODO: Add function for warnings
		}

	}
}

func addProcessor(cmd string){
	if len(cmd) > 0{
		var mPath = path.Join(core.CurrentDir, "data", "processors", cmd + ".json")
		p, err := processors.Create(mPath)
		//get instantiated module
		if err != nil {
			message("warn", err.Error())
		}	else {
			//Add instatiated module to debugger
			debugger.AddProcessor(p)
		}
	}
}

func menuConfigureProcessor(cmd string){
	if len(cmd) > 0 {
		p, errModule := interceptor.GetProcessor(cmd)
		if errModule != nil {
			message("warn", errModule.Error())
		} else {
			shellProcessor = s
			prompt.Config.AutoComplete = getCompleter("processor")
			prompt.SetPrompt("\033[31mMerlin[\033[32mmodule\033[31m][\033[33m" + shellModule.Name + "\033[31m]»\033[0m ")
			shellMenuContext = "processor"
		}
	}
}

func executeCommand(name string, arg []string) {
	var cmd *exec.Cmd

	cmd = exec.Command(name, arg...)

	out, err := cmd.CombinedOutput()

	if err != nil {
		message("warn", err.Error())
	} else {
		message("success", fmt.Sprintf("%s", out))
	}
}

func message (level string, message string) {
	switch level {
	case "info":
		color.Cyan("[i]" + message)
	case "note":
		color.Yellow("[-]" + message)
	case "warn":
		color.Red("[!]" + message)
	case "debug":
		color.Red("[DEBUG]" + message)
	case "success":
		color.Green("[+]" + message)
	default:
		color.Red("[_-_]Invalid message level: " + message)
	}
}


func menuHelpMain() {
	color.Yellow("gorp RE!!")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCaption(true, "Main Menu Help")
	table.SetHeader([]string{"Command", "Description", "Options"})

	data := [][]string{
		{"agent", "Interact with agents or list agents", "interact, list"},
		{"banner", "Print the Merlin banner", ""},
		{"exit", "Exit and close the Merlin server", ""},
		{"interact", "Interact with an agent. Alias for Empire users", ""},
		{"quit", "Exit and close the Merlin server", ""},
		{"remove", "Remove or delete a DEAD agent from the server"},
		{"sessions", "List all agents session information. Alias for MSF users", ""},
		{"use", "Use a function of Merlin", "module"},
		{"version", "Print the Merlin server version", ""},
		{"*", "Anything else will be execute on the host operating system", ""},
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}


func exit(){
	color.Red("[!]Quitting")
	serverLog.WriteString(fmt.Sprintf("[%s]Shutting down Merlin Server due to user input", time.Now()))
	os.Exit(0)
}
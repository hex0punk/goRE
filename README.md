[![Go Report Card](https://goreportcard.com/badge/github.com/DharmaOfCode/gorp)](https://goreportcard.com/report/github.com/DharmaOfCode/gorp)

# gorp
Exploring pentesting and reverse engineering uses of the Chrome DevTools protocol with Go. 

If you want to learn more about how this idea came about and how I went about writting this, you can read [this blog post](https://codedharma.com/posts/chrome-devtools-fun-with-golang/)

Right now the script intercepts requests, processes body responses and unhides hidden input. It also sets all `ng-if` and `*ngIf` attributes to true (for Angular 1.X and 2+). This can be helpful when you want to explore an Angular application and see what content is not rendered by Angular on page loads, as in some cases this could allow us to discover directories that are hidden to some users right away (note that `ng-if` does not hide or show input; rather, it decides whether to render an element on page loads). 

While there are some options (i.e. `Verbose`, `EnableConsoleLogging`, etc.), those can only be set by changing the code. This is only temporarily, as I am working on making this an actual tool and not just a script.

## Running the PoC

```
go run main.go
```
You can change the following `DebuggerOptions` flags directly in the code:

```golang
type DebuggerOptions struct {
	EnableConsole 	bool
	Verbose       	bool

	AlterDocument	bool
	AlterScript		bool
}
```
I am currentely working on a CLI tool to make use of this PoC and will be pushing updates to a feature branch for that. 

## Caveats
- The tool will crash when accessing pages when reponse bodys are ~ 10300 bytes long. Body responses this big are gzipped, and I am workong on a solution for this, but I am running into some issue. I have [question posted in StackoveFlow](https://stackoverflow.com/questions/52788269/chrome-devtools-protocol-continueinterceptedrequest-with-gzip-body-in-golang) and I keep testing ways to solve this issue. 


## Immediate Needs
- I have not found a JS beautifies and deobfuscation go library yet. Worse case scenario, I could either write one (kinda of a project of its own) or use node libraries via system calls.
- I am working on adding interactive CLI

## Todo
 
 - At the moment, this is not much more than a PoC. The idea is to make this into an actual tool that allows you do things such as:
     - Keep track of values such as user GUIDs and show alarms when certain conditions occur while you explore an application (helpful for finding IDORs).
     - Perform framework specific analysis of an application as it is explored. For instance, the tool could list all Angular services, or all API endpoints as it analyses scripts used by the application.
     - Alter scripts to test for specific behaviors, such as setting `isAdmin` variables to `true` before a request is submitted to a server.
     - Other cool stuff.
 - Add a CLI library. I am currently deciding on whether to use github.com/spf13/cobra, in which case commands may end up looking like this: `gorp intercept angular processors:ngifreplace,somemodulename` or an interactive CLI library that would allow things like this:
     ```
     >>> add modules/interceptors/angular/ngif_remove
     >>> add modules/analyzers/general/api_endpoints
     >>> start debugger
     ```
 - Add the ability for anyone to write and add custom modules. 

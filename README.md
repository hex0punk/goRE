[![Go Report Card](https://goreportcard.com/badge/github.com/DharmaOfCode/gorp)](https://goreportcard.com/report/github.com/DharmaOfCode/gorp)

# gorp
Exploring pentesting and reverse engineering uses of the Chrome DevTools protocol with Go. 

Right now the script intercepts requests, processes body responses and unhides hidden input. It also sets all `ng-if` and `*ngIf` attributes to true (for Angular 1.X and 2+). This can be helpful when you want to explore an Angular applacation and see what content is not rendered by Angular on page loads, as in some cases this could allow us to discover directories that are hidden to some users right away (note that `ng-if` does not hide or show input;rather, it decides whether to render an element on page loads). 

While there are some options (i.e. `Verbose`, `EnableConsoleLogging`, etc.), those can only be set by changing the code. This is only temporarly, as I am working on making this an actual tool and not just a script.  

## Caveats
- The tool will crash when accessing some web pages. I have not found the reason yet, though I will continue to troubleshoot it. This typically occurss here:

   ```golang
   if rawAlteredResponse != "" {
	   log.Println("[+] Sending modified body")
	   s.Debugger.ContinueInterceptedRequest(iid, godet.ErrorReason(reason), rawAlteredResponse, "", "", "", nil)
   }
   ```


## Inmediate Needs
- I have not found a JS unminification and deobfuscation go library yet. Worse case scanario, I could either write one (kinda of a project of its own) or use node libraries via system calls.
- I am working on adding interactive CLI

## Todo
 
 - At the moment the code is not much more than a PoC. The idea is to make this into an actual tool that allows you do things such as:
     - Keep track of values such as user GUIDs and show alarms when certain contions occur while you explore an application (helpful for finding IDORs).
     - Perform framework specific analysis of an application as it is explored. For instance, the tool could list all Angular services, or all API endpoints as it analyses scriptsused by the application.
     - Alter scripts to test for specific behaviors, such as setting `isAdmin` variablesto `true` before a request is submitted to a server.
     - Other cool stuff.
 - Add a CLI library. I am currentely investigating 
 - Add the ability to load custom modules. 


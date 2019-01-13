[![Go Report Card](https://goreportcard.com/badge/github.com/DharmaOfCode/gorp)](https://goreportcard.com/report/github.com/DharmaOfCode/gorp)
[![Go Documentation](http://godoc.org/github.com/DharmaOfCode/gorp?status.svg)](http://godoc.org/github.com/DharmaOfCode/gorp)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

# Gorp
A modular bug hunting, pentesting and webapp reverse engineering framework written in Go.

_If you want to learn more about how this idea came about and how I went about writing this, you can read [this blog post](https://codedharma.com/posts/chrome-devtools-fun-with-golang/). However, note that a lot has changed in the architecture since I wrote that post._

gorp is an created for web pentesting and reverse engineering. It leverages the Chrome Dev Tool protocol to intercept HTTP responses as you conduct pentest with Chrome via the use of go plugins.

## gorp plugins
Gorp plugins are essentially modules that you can use to modify or audit web responses. There are two different types of plugins (so far):

- **Processors:** processors plugins alter the response before it is rendered in the browser. This can be useful for things like modifying JavaScript code, changing HTML directives, unhiding elements in the page, highlighting areas of interest, etc.

- **Inspectors:**: inspectors conduct  analysis on responses. For instance, you may want to record all references to API calls made by the application by inspecting JavaScript code. This way, rather than waiting until the browser makes a call to `/api/admin/adduser`, you may be able to find a reference to that path in the client side code. JS Framework specific inspectors could also be used to inspect things such as services, controllers, authorization controllers, etc. Inspectors do not modify responses.


### Recompiling gorp plugins
At the moment there are constant changes on the module package. A change in that package would require that plugins are recompiled. This can be a pain as every module would need to be recompiled, so we have automated that task. Just run the below command and all modules will be recompiled:

```shell
go run main.go -p
```

## Using gorp
1. Create a configuration file that uses the structure used by the `config.yaml` file in the root directory of this repo.
2. Make sure the plugins that you want to use are compiled. You can compile all available plugins by running `go run main.go -p`
3. You can find information about any plugin by running this command:
   ```bash
   go run main.go -i -m "/the/path/of/the/module/"
   ```
4. To run gorp:
   ```bash
   go run main.g -c "./path/to/your/config/file.yml"
   ```
   
If run successfully, a new Chrome window should open up with two tabs. Use the second tab to navigate to the site that you are currently pentesting. Press `ctrl + c` to end the session (TODO: make a more effective way to end sessions).

### Ok,but what can I actually do with gorp?

There are 7 modules available at the moment. You can find information about each plugin by running `go run main.go -i /path/to/module/`. 

Here are some fun things that you can do right now. Each task is followed by a code snippet showing how your config would look like to enable the right plugins. Note that you enable multiple plugins at the same time.

**1) Force Angular 2 application to load in develoment mode**

```yaml
scope: ""
verbose: False
flags: ["-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"]
modules:
  processors:
    - path: "/data/modules/processors/angular/prodModeHijacker/"
      options: {}
```

**2) Hijack and alter a function loaded by a web application**

```yaml
scope: ""
verbose: False
flags: ["-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"]
modules:
  processors:
    - path: "/data/modules/processors/generic/functionhijacker/"
      options:
         Indicator: "isLoggedIn"
         NewBody: "return true"
```

**3) Record API calls in a file**

```yaml
scope: ""
verbose: False
flags: ["-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"]
modules:
  inspectors:
    - path: "/data/modules/inspectors/generic/apifinder/"
      options:
        FilePath : "./logs/apifinds.txt"
```

**4) Inject code in an existing function**

```yaml
scope: ""
verbose: False
flags: ["-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"]
modules:
  processors:
   - path: "/data/modules/processors/generic/injector/"
      options:
        FunctionName: "isAdmin"
        Injection: "console.log('function called, injection confirmed!');return true;"}
```

**5) Set all ngIf and ng-if attributes to always return true (applies to Angular apps)**

```
scope: ""
verbose: False
flags: ["-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"]
modules:
  processors:
    - path: "/data/modules/processors/angular/unhider/"
      options: {}
```


**6) Simple find and replace**
```
scope: ""
verbose: False
flags: ["-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"]
modules:
  processors:
    - path: "/data/modules/processors/generic/findreplace/"
      options:
         Find: "isAdmin=false"
         Replace: "isAdmin=true"
```

**7) Unhide all hidden input and add highlight what the input is used for
```
scope: ""
verbose: False
flags: ["-na", "--disable-gpu", "--window-size=1200,800", "--auto-open-devtools-for-tabs","--disable-popup-blocking"]
modules:
  processors:
    - path: "/data/modules/processors/generic/unhider/"
      options: {}
```

## Creating your own gorp plugin
The power of gorp is in the plugins. Creating your own plugin is simple.

1. Create a file called `gorpmod.go` under `/data/modules/processors` or `/data/modules/inspectors`, depending on your type of plugin (see above for the differences between an inspector and a processor.
2. Depending on the type of plugin, your code must implement either the `Processor` or `Inspector` interface, which are declared in the `modules` package. Both module types must accept a struct parameter of type `modules.WebData` which gives your module access the response body, headers and type. The type can be `Document`, `Script` or `Request` (`Request` types have not been implemented yet but that is my list of priorities for this gorp).
3. Your plugin must include a symbol to be used by gorp. The symbol should be declared like this:

   ```golang
   //apifinder is just the name of your plugin
   type apifinder struct {
       Registry    modules.Registry
       Options    []modules.Option
   }
   ```
4. Make sure to export the symbol at the end of your plugin, like so:

   ```golang
   var Inspector apifinder
   ```
 5. Compile your plugin like so:
 
    ```bash
    go build -buildmode=plugin -o gorpmod.so gorpmod.go
    ```
 6. Now you are ready to use your plugin with gorp. 


## Immediate Needs
- I have not found a JS beautifies and deobfuscation go library yet. Worst-case scenario, I could either write one (kinda of a project of its own) or use node libraries via system calls.

## Todo
 
 - Add a fancy, interactive shell-like CLI. 
 - Rad CLI colors and functions for fancy cli printing
 - Create more plugins for tasks such as:
     - Keep track of values such as user GUIDs and show alarms when certain conditions occur while you explore an application (helpful for finding IDORs).
     - Perform framework specific analysis of an application as it is explored. For instance, the tool could list all Angular services or all API endpoints as it analyses scripts used by the application.
     - Other rad stuff.

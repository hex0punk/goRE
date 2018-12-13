[![Go Report Card](https://goreportcard.com/badge/github.com/DharmaOfCode/gorp)](https://goreportcard.com/report/github.com/DharmaOfCode/gorp)

# gorp
AppSec pentesting and reverse engineering tool that uses the Chrome DevTools protocol with Go. 

_If you want to learn more about how this idea came about and how I went about writing this, you can read [this blog post](https://codedharma.com/posts/chrome-devtools-fun-with-golang/). However, note that a lot has changed in the architecture since I wrote that post.

gorp is a Chrome dev protocol engine created for pentesters and hackers. It leverages the Chrome Dev Tool protocol to intercept HTTP responses as you conduct pentest with Chrome via the use of go plugins.

## gorp plugin
Gorp plugins are essentially modules that you can use to modify or audit web responses. There are two different types of plugins (so far):

- **Processors:** processors plugins alter the response before it is rendered in the browser. This can be useful for things like modifying JavaScript code, changing HTML directives, unhiding elements in the page, highlighting areas of interest, etc.

- **Inspectors:**: inspectors conduct  analysis on responses. For instance, you may want to record all references to API calls made by the application by inspecting JavaScript code. This way, rather than waiting until the browser makes a call to `/api/admin/adduser`, you may be able to find a reference to that path in the client side code. JS Framework specific inspectors could also be used to inspect things such as services, controllers, authorization controllers, etc. Inspectors do not modify responses.

## Using gorp
1. Create a configuration file that uses the structure used by the `config.yaml` file in the root directory of this repo.
2. You can find information about any plugin by running this command:
   ```bash
   go run main.go -i -m "/the/path/of/the/module/"
   ```
3. To run gorp:
   ```bash
   go run main.g -c "./path/to/your/config/file.yml"
   ```
   
If run successfully, a new Chrome window should open up with two tabs. Use the second tab to navigate to the site that you are currently pentesting. Press `ctrl + c` to end the session (TODO: make a more effective way to end sessions).

## Creating your own gorp plugin
The power of gorp is in the plugins. Creating your own plugin is simple.

1. Create a file under `` or ``, depending on your type of plugin (see above for the differences between an inspector and a processor.
2. Depending on the type of plugin, your code must implement either the `Processor` or `Inspector` interface, which are declared in the `modules` package.
3. Your plugin must return a symbol to be used by gorp. The symbol should be declared like this:

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
 
 - Go doc
 - Create a new type of module that works on requests. Currently processors and inspectors work only on responses.
 - Add a fancy, interactive shell-like CLI. 
 - Create more plugins for tasks such as:
     - Keep track of values such as user GUIDs and show alarms when certain conditions occur while you explore an application (helpful for finding IDORs).
     - Perform framework specific analysis of an application as it is explored. For instance, the tool could list all Angular services or all API endpoints as it analyses scripts used by the application.
     - Alter scripts to test for specific behaviors, such as setting `isAdmin` variables to `true` before a request is submitted to a server.
     - Other cool stuff.

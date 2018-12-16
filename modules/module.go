// Package modules provides functionality for for loading, configuring and running gorp modules
package modules

import (
	"fmt"
	"github.com/DharmaOfCode/gorp/base"
	"github.com/fatih/color"
	"plugin"
)

// Modules holds selected processors and inspectors to be used in a gorp session
type Modules struct {
	Processors	[]ProcessorModule
	Inspectors  []InspectorModule
}

// Registry holds meta data for modules
type Registry struct {
	Name     string  	`json:"name"` 	// Name of the module
	DocTypes []string 	`json:"docTypes"` 	// Type of document
	Author   []string 	`json:"author"`	// A list of module authors
	Credits	 []string	`json:"credits"` // A list of people to credit for underlying tool or techniques
	Path     string 	`json:"path"`	// Path to the module
	Description string 	`json:"description"`	// A description of what the module does
	Notes    string 	`json:"notes"`	// Additional information or notes about the module
}

// Option contains options specific to modules
type Option struct {
	Name 		string		`json:"name"` 		// Name of the option
	Value 		string		`json:"value"` 		// Value of the option
	Required 	bool		`json:"required"` 	// Is this a required option?
	Description string		`json:"description"`// A description of the option
}

// ProcessorModule represents a processor module. Processor modules alter the body of a request or response
type ProcessorModule struct {
	Process	func(webData WebData) (string, error)
	Registry	Registry
	Options 	[]Option 	`json:"options"`	// A list of configurable options/arguments for the module
}

// InspectorModule represents an inspector module. Inspectors analyse responses to answer questions about the
// application or to discover different types of information found in HTML documents, JavaScript comments and code
type InspectorModule struct {
	Inspect	func(webData WebData) error
	Registry	Registry
	Options 	[]Option
}

// Processor identifies the functions that all processor modules must implement.
type Processor interface {
	Init() // Init Initializes module data
	GetOptions()  []Option // GetOptions returns a list of available options for the module
	GetRegistry() Registry // GetRegistry returns an object with meta data describing the module
	Process(webData WebData) (string, error) // Process alters the body of a request
}

// Inspector identifies the functions that all inspector modules must implement.
type Inspector interface {
	Init() // Init Initializes module data
	GetOptions()  []Option  // GetOptions returns a list of available options for the module
	GetRegistry() Registry // GetRegistry returns an object with meta data describing the module
	Inspect(webData WebData) error // Inspect inspects web content for discovery and recon purposes
}

// WebData identifies a web request or response object. The type can be either "Document," "Script," or "Request"
type WebData struct {
	Body	string
	Headers map[string]interface{}
	Type	string
	Url		string
}

// InitProcessors initializes modules selected for a gorp session
func (m *Modules) InitProcessors(mods []base.ModuleConfig) error{
	for _, v := range mods{
		module, err := m.GetProcessor(v.Path)
		if err != nil {
			return err
		}

		for option, value := range v.Options{
			_, err := module.SetOption(option, value)
			if err != nil {
				return err
			}
		}
		m.Processors = append(m.Processors, *module)
	}
	return nil
}

// GetProcessor looks up and loads a processor module as Go plugins.
// It returns a pointer to the processor module
func (m *Modules) GetProcessor(path string) (*ProcessorModule, error) {
	module := ProcessorModule{}
	fmt.Println("[+] Loading module: " + path)
	mod := "." + path + "gorpmod.so"
	plug, err := plugin.Open(mod)
	if err != nil {
		return nil, err
	}

	// look up a symbol (an exported function or variable)
	// in this case, variable Greeter
	symProcessor, err := plug.Lookup("Processor")
	if err != nil {
		return nil, err
	}

	// Assert that loaded symbol is of a desired type
	// in this case interface type Greeter (defined above)
	var processor Processor
	//processor = new(modules.Processor)
	processor, ok := symProcessor.(Processor)
	if !ok {
		fmt.Println("unexpected type from processor symbol")
		return nil, err
	}
	processor.Init()
	module.Registry = processor.GetRegistry()
	module.Options = processor.GetOptions()
	module.Process = processor.Process
	return &module, nil
}

// InitInspectors  loads a list of inspector modules.
func (m *Modules) InitInspectors(mods []base.ModuleConfig) error {
	for _, v := range mods{
		module, err := m.GetInspector(v.Path)
		if err != nil {
			return err
		}

		for option, value := range v.Options{
			_, err := module.SetOption(option, value)
			if err != nil {
				return err
			}
		}
		m.Inspectors = append(m.Inspectors, *module)
	}
	return nil
}

// GetInspector looks up and loads an inspector module as Go plugins.
// It returns a pointer to the inspector module
func (m *Modules) GetInspector(path string) (*InspectorModule, error) {
		module := InspectorModule{}
		fmt.Println("[+] Loading module: " + path)
		mod := "." + path + "gorpmod.so"
		plug, err := plugin.Open(mod)
		if err != nil {
			return nil, err
		}

		// look up a symbol (an exported function or variable)
		// in this case, variable Greeter
		symProcessor, err := plug.Lookup("Inspector")
		if err != nil {
			return nil, err
		}

		// Assert that loaded symbol is of a desired type
		// in this case interface type Greeter (defined above)
		var inspector Inspector
		//processor = new(modules.Processor)
		inspector, ok := symProcessor.(Inspector)
		if !ok {
			fmt.Println("unexpected type from processor symbol")
			return nil, err
		}
		inspector.Init()
		module.Registry = inspector.GetRegistry()
		module.Options = inspector.GetOptions()
		module.Inspect = inspector.Inspect
		return &module, nil
}

// ShowInfo displays the information for the given processor module
func (p *ProcessorModule) ShowInfo(){
	showInfo(p.Registry)
}

// ShowInfo displays the information for the given inspector module
func (i *InspectorModule) ShowInfo(){
	showInfo(i.Registry)
}

// SetOption is used to change and set a processor module option. Used when a user is configuring a processor module.
// It returns a message indicating whether the options was set successfully.
func (p *ProcessorModule) SetOption(option string, value string) (string, error){
	// Verify this option exists
	for k, v := range p.Options {
		if option == v.Name {
			p.Options[k].Value = value
			return fmt.Sprintf("%s set to %s", v.Name, p.Options[k].Value), nil
		}
	}
	return "", fmt.Errorf("invalid module option: %s", option)
}

// SetOption is used to change and set an inspector module option. Used when a user is configuring an inspector module.
// It returns a message indicating whether the options was set successfully.
func (i *InspectorModule) SetOption(option string, value string) (string, error){
	// Verify this option exists
	for k, v := range i.Options {
		if option == v.Name {
			i.Options[k].Value = value
			return fmt.Sprintf("%s set to %s", v.Name, i.Options[k].Value), nil
		}
	}
	return "", fmt.Errorf("invalid module option: %s", option)
}

func showInfo(r Registry){
	color.Yellow("Module:\r\n\t%s\r\n", r.Name)
	color.Yellow("Authors:")
	for d := range r.DocTypes{
		color.Yellow("\t%s", r.DocTypes[d])
	}
	color.Yellow("Authors:")
	for a := range r.Author {
		color.Yellow("\t%s", r.Author[a])
	}
	color.Yellow("Credits:")
	for c := range r.Credits {
		color.Yellow("\t%s", r.Credits[c])
	}
	color.Yellow("Description:\r\n\t%s", r.Description)
	fmt.Println()
	color.Yellow("Notes: %s", r.Notes)
}

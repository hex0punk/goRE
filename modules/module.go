package modules

import (
	"fmt"
	"github.com/fatih/color"
	"plugin"
)

type Modules struct {
	Processors	[]ProcessorModule
	Inspectors  []InspectorModule
}

// Module is a structure containing the base information or template for modules
type Registry struct {
	Name     string  	`json:"name"` 	// Name of the module
	Author   []string 	`json:"author"`	// A list of module authors
	Credits	 []string	`json:"credits"` // A list of people to credit for underlying tool or techniques
	Path     string 	`json:"path"`	// Path to the module
	Description string 	`json:"description"`	// A description of what the module does
	Notes    string 	`json:"notes"`	// Additional information or notes about the module
}

// Option is a structure containing the keys for the object
type Option struct {
	Name 		string		`json:"name"` 		// Name of the option
	Value 		string		`json:"value"` 		// Value of the option
	Required 	bool		`json:"required"` 	// Is this a required option?
	Flag 		string		`json:"flag"`		// The command line flag used for the option
	Description string		`json:"description"`// A description of the option
}
type ProcessorModule struct {
	Process	func(body string) (string, error)
	Registry	Registry
	Options 	[]Option 	`json:"options"`	// A list of configurable options/arguments for the module
}

type InspectorModule struct {
	Inspect	func(body string) error
	Registry	Registry
	Options 	[]Option
}

type Processor interface {
	Init()
	GetOptions()  []Option
	GetRegistry() Registry
	Process(body string) (string, error)
}

type Inspector interface {
	Init()
	GetOptions()  []Option
	GetRegistry() Registry
	Inspect(body string) error
}

func (m *Modules) InitProcessors(paths []string) error{
	for _, v := range paths{
		module := ProcessorModule{}
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
		var processor Processor
		//processor = new(modules.Processor)
		processor, ok := symProcessor.(Processor)
		if !ok {
			fmt.Println("unexpected type from processor symbol")
			return err
		}
		processor.Init()
		module.Registry = processor.GetRegistry()
		module.Process = processor.Process
		m.Processors = append(m.Processors, module)
	}
	return nil
}

func (m *Modules) InitInspectors(paths []string) error {
	for _, v := range paths{
		module := InspectorModule{}
		fmt.Println("[+] Loading module: " + v)
		mod := v + "gorpmod.so"
		plug, err := plugin.Open(mod)
		if err != nil {
			return err
		}

		// look up a symbol (an exported function or variable)
		// in this case, variable Greeter
		symProcessor, err := plug.Lookup("Inspector")
		if err != nil {
			return err
		}

		// Assert that loaded symbol is of a desired type
		// in this case interface type Greeter (defined above)
		var inspector Inspector
		//processor = new(modules.Processor)
		inspector, ok := symProcessor.(Inspector)
		if !ok {
			fmt.Println("unexpected type from processor symbol")
			return err
		}
		inspector.Init()
		module.Registry = inspector.GetRegistry()
		module.Inspect = inspector.Inspect
		m.Inspectors = append(m.Inspectors, module)
	}

	return nil
}

func (p *ProcessorModule) ShowInfo(){
	showInfo(p.Registry)
}

func (i *InspectorModule) ShowInfo(){
	showInfo(i.Registry)
}

// SetOption is used to change the passed in module option's value. Used when a user is configuring a module
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

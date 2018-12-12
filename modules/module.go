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

type ProcessorModule struct {
	Process	func(body string) (string, error)
	Registry	Registry
}

type InspectorModule struct {
	Inspector	Inspector
	Registry	Registry
}

type Processor interface {
	GetRegistry() Registry
	Process(body string) (string, error)
}

type Inspector interface {
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
		module.Registry = inspector.GetRegistry()
		module.Inspector = inspector
		m.Inspectors = append(m.Inspectors, module)
	}

	return nil
}

func (p *ProcessorModule) ShowInfo(){
	showInfo(p.Registry)
}

func (p *InspectorModule) ShowInfo(){
	showInfo(p.Registry)
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

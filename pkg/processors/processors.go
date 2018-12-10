package processors

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Processor struct {
	Name     string  	`json:"name"`
	Author   []string 	`json:"author"`
	Credits	 []string	`json:"credits"`
	Path     []string 	`json:"path"`
	Description string 	`json:"description"`
	Commands []string 	`json:"commands"`
	SourceLocal	[]string 	`json:"local"`
	Options []Option 	`json:"options"`
}

type Option struct {
	Name 		string		`json:"name"`
	Value 		string		`json:"value"`
	Required 	bool		`json:"required"`
	Flag 		string		`json:"flag"`
	Description string		`json:"description"`
}


func Create(modulePath string) (Processor, error) {
	var p Processor

	f, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return m, err
	}

	// Unmarshal processor's JSON message
	var moduleJSON map[string]*json.RawMessage
	errModule := json.Unmarshal(f, &moduleJSON)
	if errModule != nil {
		return p, errModule
	}

	// Determine all message types
	var keys []string
	for k := range moduleJSON {
		keys = append(keys,k)
	}

	// Validate that procesor's JSON contains at least the base message
	var containsBase bool
	for i := range keys{
		if keys[i] == "base" {
			containsBase = true
		}
	}

	// Marshal Base message type
	if !containsBase {
		return p, errors.New("the module's definition does not contain the 'BASE' message type")
	}
	errJSON := json.Unmarshal(*moduleJSON["base"], &p)
	if errJSON != nil {
		return p, errJSON
	}

	return p, nil
}

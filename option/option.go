package option

import "strings"

// Option contains options specific to modules
type Option struct {
	Name        string `json:"name"`        // Name of the option
	Value       string `json:"value"`       // Value of the option
	Required    bool   `json:"required"`    // Is this a required option?
	Description string `json:"description"` // A description of the option
}


func (o *Option) IsList() bool {
	return strings.Contains(o.Value, ",")
}

func (o *Option) GetAsList() []string{
	if o.IsList(){
		return strings.Split(o.Value, ",")
	} else {
		return nil
	}
}
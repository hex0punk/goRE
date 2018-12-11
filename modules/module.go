package modules

type Processor interface {
	Process(body string) (string, error)
}

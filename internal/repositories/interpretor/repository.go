package interpretor

type Interpretor interface {
	Interpret(text string) (InterpretorResponse, error)
}

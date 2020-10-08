package interpreter

type Response struct {
	ID       string
	Text     string
	Intents  []Intent
	Entities map[string]interface{}
	Traits   map[string]interface{}
}

type Intent struct {
	Name       string
	Confidence float64
}

type Interpreter interface {
	InterpretText(text string) (Response, error)
}

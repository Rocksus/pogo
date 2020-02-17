package interpretor

import witai "github.com/wit-ai/wit-go"

type InterpretorResponse struct {
	ID         string
	Text       string
	Intent     string
	Confidence float64
	Entities   map[string]interface{}
}

type interpretor struct {
	client *witai.Client
}

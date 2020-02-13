package interpretor

type InterpretorResponse struct {
	Intent     string
	Confidence float64
	Entities   []string
}

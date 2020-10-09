package witai

import (
	"github.com/Rocksus/pogo/configs"
	"github.com/Rocksus/pogo/internal/repositories/interpreter"
	"github.com/mitchellh/mapstructure"
	witai "github.com/wit-ai/wit-go"
)

type witInterpreter struct {
	client *witai.Client
}

func NewInterpreter(config configs.InterpretorConfig) interpreter.Interpreter {
	client := witai.NewClient(config.Token)
	return &witInterpreter{
		client: client,
	}
}

func (w *witInterpreter) InterpretText(text string) (interpreter.Response, error) {
	msg, err := w.client.Parse(&witai.MessageRequest{
		Query: text,
	})
	if err != nil {
		return interpreter.Response{}, err
	}

	return w.parseResponse(msg)
}

func (w *witInterpreter) parseResponse(resp *witai.MessageResponse) (res interpreter.Response, err error) {
	// TODO: api version is old. This parser might need to change later on
	type intent struct {
		Confidence float64 `mapstructure:"confidence"`
		Value      string  `mapstructure:"value"`
	}

	type entities struct {
		Intents []intent               `mapstructure:"intent"`
		Extra   map[string]interface{} `mapstructure:",remain"`
	}

	var e entities
	err = mapstructure.Decode(resp.Entities, &e)
	if err != nil {
		return
	}

	res = interpreter.Response{
		ID:       resp.ID,
		Text:     resp.Text,
		Entities: e.Extra,
	}
	for _, i := range e.Intents {
		res.Intents = append(res.Intents, interpreter.Intent{Name: i.Value, Confidence: i.Confidence})
	}
	return
}

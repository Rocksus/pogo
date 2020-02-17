package interpretor

import (
	"fmt"

	"github.com/Rocksus/pogo/configs"
	witai "github.com/wit-ai/wit-go"
)

type Interpretor interface {
	InterpretText(text string) (*InterpretorResponse, error)
}

func InitInterpretorRepository(config configs.InterpretorConfig) Interpretor {
	client := witai.NewClient(config.Token)
	return &interpretor{
		client: client,
	}
}

func (i *interpretor) InterpretText(text string) (*InterpretorResponse, error) {
	msg, err := i.client.Parse(&witai.MessageRequest{
		Query: text,
	})
	if err != nil {
		return nil, err
	}
	resp := &InterpretorResponse{
		ID:       msg.ID,
		Text:     msg.Text,
		Entities: msg.Entities,
	}
	fmt.Println(msg.Entities)

	return resp, nil
}

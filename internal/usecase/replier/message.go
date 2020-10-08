package replier

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/Rocksus/pogo/internal/modules/joke"
	"github.com/Rocksus/pogo/internal/modules/weather"
	"github.com/Rocksus/pogo/internal/repositories/interpreter"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
)

type MessageReplier interface {
	Reply(ctx context.Context, message linebot.Message) (reply linebot.SendingMessage)
}

type messageReplier struct {
	interpreter interpreter.Interpreter
}

func NewMessageReplier(interpreter interpreter.Interpreter) MessageReplier {
	return &messageReplier{
		interpreter: interpreter,
	}
}

func (m *messageReplier) Reply(ctx context.Context, message linebot.Message) (reply linebot.SendingMessage) {
	var err error
	switch msg := message.(type) {
	case *linebot.TextMessage:
		reply, err = m.handleTextMessage(ctx, msg)
	case *linebot.StickerMessage:
		reply, err = m.handleStickerMessage(ctx, msg)
	default:
		err = fmt.Errorf("unhandled message type: %s", m.getMessageType(message))
		log.Errorln(err)
	}
	if err != nil {
		return m.createDefaultReply()
	}

	return reply
}

func (m *messageReplier) handleTextMessage(ctx context.Context, msg *linebot.TextMessage) (reply linebot.SendingMessage, err error) {
	data, err := m.interpreter.InterpretText(msg.Text)
	if err != nil {
		log.Errorln(err)
		return
	}

	intent := m.getBestIntent(data)
	// TODO: use map of handlers later. This is just a temporary logic
	switch intent.Name {
	case "tellJoke":
		j, err := joke.GetRandomJoke()
		if err != nil {
			return nil, err
		}
		jokeStr := fmt.Sprintf("%s\n\n%s", j.Setup, j.Punchline)
		return linebot.NewTextMessage(jokeStr), nil
	case "weather/checkWeather":
		w, err := weather.QueryLocation("jakarta")
		if err != nil {
			return nil, err
		}
		weatherStr := fmt.Sprintf(
			"Here's the weather in %s, %s:\n\n%s\nHumidity: %d\nTemperature: %.2f°C",
			w.Name,
			w.System.Country,
			w.Weather[0].Description,
			w.Details.Humidity,
			w.Details.TemperatureCelcius,
		)
		return linebot.NewTextMessage(weatherStr), nil
	default:
		return nil, fmt.Errorf("unhandled intent: %s", intent.Name)
	}
}

func (m *messageReplier) handleStickerMessage(ctx context.Context, msg *linebot.StickerMessage) (reply linebot.SendingMessage, err error) {
	// This is just to mimic the previous behavior
	replyText := fmt.Sprintf(
		"sticker id is %s, stickerResourceType is %s",
		msg.StickerID,
		msg.StickerResourceType,
	)
	return linebot.NewTextMessage(replyText), nil
}

func (m *messageReplier) getBestIntent(resp interpreter.Response) interpreter.Intent {
	if len(resp.Intents) == 0 {
		return interpreter.Intent{}
	}

	intents := make([]interpreter.Intent, len(resp.Intents))
	copy(intents, resp.Intents)
	sort.Slice(intents, func(i, j int) bool {
		return intents[i].Confidence > intents[j].Confidence
	})

	return intents[0]
}

func (m *messageReplier) createDefaultReply() linebot.SendingMessage {
	return linebot.NewTextMessage("Sorry, I don't quite get that.")
}

func (m *messageReplier) getMessageType(msg linebot.Message) string {
	t := reflect.TypeOf(msg)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

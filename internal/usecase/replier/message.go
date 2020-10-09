package replier

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/Rocksus/pogo/internal/repositories/interpreter"
	"github.com/Rocksus/pogo/pkg/plugin"
	"github.com/Rocksus/pogo/pkg/plugin/weather"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
)

type MessageReplier interface {
	Reply(ctx context.Context, event *linebot.Event) chan linebot.SendingMessage
}

type messageReplier struct {
	interpreter interpreter.Interpreter
	plugins     map[string]plugin.MessageReplier
}

func NewMessageReplier(interpreter interpreter.Interpreter, plugins map[string]plugin.MessageReplier) MessageReplier {
	return &messageReplier{
		interpreter: interpreter,
		plugins:     plugins,
	}
}

func (m *messageReplier) Reply(ctx context.Context, event *linebot.Event) (replyCh chan linebot.SendingMessage) {
	// safeguard against non-message events
	if event.Type != linebot.EventTypeMessage {
		return
	}
	message := event.Message

	replyCh = make(chan linebot.SendingMessage)
	go func() {
		defer close(replyCh)

		switch msg := message.(type) {
		case *linebot.TextMessage:
			m.handleTextMessage(ctx, msg, replyCh)
		case *linebot.StickerMessage:
			m.handleStickerMessage(ctx, msg, replyCh)
		default:
			log.Errorln("unhandled message type:", m.getMessageType(msg))
			replyCh <- m.createDefaultReply()
		}
	}()

	return replyCh
}

func (m *messageReplier) handleTextMessage(ctx context.Context, msg *linebot.TextMessage, replyCh chan<- linebot.SendingMessage) {
	data, err := m.interpreter.InterpretText(msg.Text)
	if err != nil {
		log.Errorln(err)
		return
	}

	intent := m.getBestIntent(data)
	replier, ok := m.plugins[intent.Name]
	if ok {
		replier.Reply(ctx, plugin.Message{
			Text:     msg.Text,
			Intent:   intent.Name,
			Entities: data.Entities,
		}, replyCh)
		return
	}

	// TODO: use map of handlers later. This is just a temporary logic
	switch intent.Name {
	case "weather/checkWeather":
		replyCh <- linebot.NewTextMessage("Hold on, let me ask the weather gods")
		time.Sleep(2 * time.Second)
		w, err := weather.QueryLocation("jakarta")
		if err != nil {
			log.Errorln(err)
			replyCh <- linebot.NewTextMessage("Sorry, the weather gods aren't answering my questions.")
			return
		}
		replyCh <- linebot.NewTextMessage(fmt.Sprintf("Got it! Here's the weather in %s, %s", w.Name, w.System.Country))
		replyCh <- linebot.NewTextMessage(fmt.Sprintf("%s: %s", w.Weather[0].Main, w.Weather[0].Description))
		replyCh <- linebot.NewTextMessage(fmt.Sprintf("Humidity: %d", w.Details.Humidity))
		replyCh <- linebot.NewTextMessage(fmt.Sprintf("Temperature: %.2f°C", w.Details.TemperatureCelcius))
	default:
		log.Errorln("unhandled intent:", intent.Name)
		replyCh <- m.createDefaultReply()
	}

	return
}

func (m *messageReplier) handleStickerMessage(ctx context.Context, msg *linebot.StickerMessage, replyCh chan<- linebot.SendingMessage) {
	// This is just to mimic the previous behavior
	replyText := fmt.Sprintf(
		"sticker id is %s, stickerResourceType is %s",
		msg.StickerID,
		msg.StickerResourceType,
	)
	replyCh <- linebot.NewTextMessage(replyText)
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

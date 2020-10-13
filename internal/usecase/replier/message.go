package replier

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/Rocksus/pogo/internal/repositories/interpreter"
	"github.com/Rocksus/pogo/pkg/plugin"
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
	if !ok {
		log.Errorln("unhandled intent:", intent.Name)
		replyCh <- m.createDefaultReply()
		return
	}

	replier.Reply(ctx, plugin.Message{
		Text:     msg.Text,
		Intent:   intent.Name,
		Entities: data.Entities,
	}, replyCh)
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

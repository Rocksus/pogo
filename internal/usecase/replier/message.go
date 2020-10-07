package replier

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Rocksus/pogo/internal/repositories/interpretor"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
)

type MessageReplier interface {
	Reply(ctx context.Context, message linebot.Message) (reply linebot.SendingMessage)
}

type messageReplier struct {
	interpreter interpretor.Interpretor
}

func NewMessageReplier(interpreter interpretor.Interpretor) MessageReplier {
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
		err = fmt.Errorf("unhandled message type: %s", getMessageType(message))
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

	// TODO: do something with data later
	_ = data

	// For now, just echo back the incoming message
	reply = linebot.NewTextMessage(msg.Text)
	return
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

func (m *messageReplier) createDefaultReply() linebot.SendingMessage {
	return linebot.NewTextMessage("Sorry, I don't quite get that.")
}

func getMessageType(msg linebot.Message) string {
	t := reflect.TypeOf(msg)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

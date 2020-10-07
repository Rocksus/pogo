package linehttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rocksus/pogo/internal/repositories/interpretor"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
)

type controller struct {
	client      *linebot.Client
	interpreter interpretor.Interpretor
}

// NewController is...
//
// TODO: intrepreter should be in use case.
func NewController(client *linebot.Client, interpreter interpretor.Interpretor) *controller {
	return &controller{
		client:      client,
		interpreter: interpreter,
	}
}

func (c *controller) HandleWebhook(w http.ResponseWriter, req *http.Request) {
	events, err := c.client.ParseRequest(req)
	if err != nil {
		log.Errorln(err)
		switch err {
		case linebot.ErrInvalidSignature:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	for _, e := range events {
		switch e.Type {
		case linebot.EventTypeMessage:
			c.handleMessageEvent(req.Context(), e)
		default:
		}
	}
}

// TODO: this function should be refactored to use the UC.
func (c *controller) handleMessageEvent(ctx context.Context, event *linebot.Event) {
	switch msg := event.Message.(type) {
	case *linebot.TextMessage:
		data, err := c.interpreter.InterpretText(msg.Text)
		if err != nil {
			log.Errorln(err)
			c.replyDefaultMessage(event.ReplyToken)
		}
		switch data.Intent {

		}

		_, err = c.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg.Text)).
			WithContext(ctx).
			Do()
		if err != nil {
			log.Errorln(err)
		}
	case *linebot.StickerMessage:
		replyMessage := fmt.Sprintf(
			"sticker id is %s, stickerResourceType is %s", msg.StickerID, msg.StickerResourceType)
		if _, err := c.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			log.Errorln(err)
		}
	}
}

func (c *controller) replyDefaultMessage(replyToken string) {
	if _, err := c.client.ReplyMessage(replyToken, linebot.NewTextMessage("Sorry I don't quite get that.")).Do(); err != nil {
		log.Errorln(err)
	}
}

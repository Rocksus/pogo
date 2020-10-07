package linehttp

import (
	"context"
	"net/http"

	"github.com/Rocksus/pogo/internal/usecase/replier"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
)

type Controller struct {
	client     *linebot.Client
	msgReplier replier.MessageReplier
}

// NewController is...
func NewController(client *linebot.Client, msgReplier replier.MessageReplier) *Controller {
	return &Controller{
		client:     client,
		msgReplier: msgReplier,
	}
}

func (c *Controller) HandleWebhook(w http.ResponseWriter, req *http.Request) {
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

func (c *Controller) handleMessageEvent(ctx context.Context, event *linebot.Event) {
	reply := c.msgReplier.Reply(ctx, event.Message)
	_, err := c.client.ReplyMessage(event.ReplyToken, reply).Do()
	if err != nil {
		log.Errorln(err)
	}
}

func (c *Controller) replyDefaultMessage(replyToken string) {
	if _, err := c.client.ReplyMessage(replyToken, linebot.NewTextMessage("Sorry I don't quite get that.")).Do(); err != nil {
		log.Errorln(err)
	}
}

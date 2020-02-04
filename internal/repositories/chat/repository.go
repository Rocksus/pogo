package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Rocksus/pogo/configs"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Repository interface {
	InitClient() error
	GetHandler(ctx context.Context) func(w http.ResponseWriter, req *http.Request)
}

func InitChatRepository(config configs.ChatConfig) Repository {
	return &lineRepo{
		MasterID:           config.MasterID,
		ChannelAccessToken: config.ChannelAccessToken,
		ChannelSecret:      config.ChannelSecret,
	}
}

func (l *lineRepo) InitClient() error {
	if l.Client != nil {
		return nil
	}
	bot, err := linebot.New(l.ChannelSecret, l.ChannelAccessToken)
	if err != nil {
		return err
	}
	l.Client = bot
	return nil
}

func (l *lineRepo) GetHandler(ctx context.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		events, err := l.Client.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = l.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = l.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
package chat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rocksus/pogo/configs"
	"github.com/Rocksus/pogo/internal/repositories/interpretor"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Repository interface {
	initClient() error
	GetHandler() func(w http.ResponseWriter, req *http.Request)
	replyDefaultMessage(replyToken string)
}

func InitChatRepository(config configs.ChatConfig, interpretor interpretor.Interpretor) Repository {
	newRepo := &lineRepo{
		MasterID:           config.MasterID,
		ChannelAccessToken: config.ChannelAccessToken,
		ChannelSecret:      config.ChannelSecret,
		Interpretor:        interpretor,
	}
	err := newRepo.initClient()
	if err != nil {
		log.Fatalf("[Init Chat] Failed to initialize chat repository, err: %s", err.Error())
	}
	log.Print("[Init Chat] Successfully initialized chat repository.")
	return newRepo
}

func (l *lineRepo) initClient() error {
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

func (l *lineRepo) replyDefaultMessage(replyToken string) {
	if _, err := l.Client.ReplyMessage(replyToken, linebot.NewTextMessage("Sorry I don't quite get that.")).Do(); err != nil {
		log.Print(err)
	}
}

func (l *lineRepo) GetHandler() func(w http.ResponseWriter, req *http.Request) {
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
					data, err := l.Interpretor.InterpretText(message.Text)
					if err != nil {
						log.Print(err)
						l.replyDefaultMessage(event.ReplyToken)
					}
					switch data.Intent {

					}
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

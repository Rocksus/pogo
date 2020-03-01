package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Rocksus/pogo/internal/modules/joke"
	"github.com/Rocksus/pogo/internal/modules/news"
	"github.com/Rocksus/pogo/internal/modules/weather"

	"github.com/Rocksus/pogo/configs"
	"github.com/Rocksus/pogo/internal/repositories/interpretor"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Repository interface {
	initClient() error
	GetHandler() func(w http.ResponseWriter, req *http.Request)
	SendDailyMessage()
	replyDefaultMessage(replyToken string)
	GetUserProfile(userID string) (*UserData, error)
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

					rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
					defer cancel()
					if _, err = l.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).WithContext(rctx).Do(); err != nil {
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

func (l *lineRepo) GetUserProfile(userID string) (*UserData, error) {
	if userID == "" {
		return nil, fmt.Errorf("[Chat][GetUserProfile] UserID can't be empty")
	}

	rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	userData, err := l.Client.GetProfile(userID).WithContext(rctx).Do()
	if err != nil {
		return nil, fmt.Errorf("[Chat][GetUserProfile] Can't get user profile, err: %s", err.Error())
	}

	return &UserData{
		UserID:        userData.UserID,
		DisplayName:   userData.DisplayName,
		PictureURL:    userData.PictureURL,
		StatusMessage: userData.StatusMessage,
	}, nil
}

func (l *lineRepo) SendDailyMessage() {
	var messages []linebot.SendingMessage
	messages = make([]linebot.SendingMessage, 3)

	jokeData, err := joke.GetRandomJoke()
	if err != nil {
		log.Printf("[Chat Cron][SendDailyMessage] Failed to get joke data, err: %s", err.Error())
	}
	weatherData, err := weather.QueryLocation("jakarta")
	if err != nil {
		log.Printf("[Chat Cron][SendDailyMessage] Failed to get weather data, err: %s", err.Error())
	}
	newsData, err := news.GetTopNews(news.TopNewsRequestParam{})
	if err != nil {
		log.Printf("[Chat Cron][SendDailyMessage] Failed to get news data, err: %s", err.Error())
	}

	messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("Hello %s, here are your daily stuffs")))
	messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("Weather data")))
	messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("%s\n\n%s", jokeData.Setup, jokeData.Punchline)))

	_, err = l.Client.PushMessage(l.MasterID, messages...).Do()
	if err != nil {
		log.Printf("[Chat Cron][SendDailyMessage] Failed to send message, err: %s", err.Error())
	}
}

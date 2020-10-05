package chat

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Rocksus/pogo/internal/modules/joke"
	"github.com/Rocksus/pogo/internal/modules/news"
	"github.com/Rocksus/pogo/internal/modules/weather"
	"github.com/Rocksus/pogo/internal/utils/stringformat"
	"github.com/nickylogan/go-log"

	"github.com/Rocksus/pogo/configs"
	"github.com/Rocksus/pogo/internal/repositories/interpretor"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Repository interface {
	GetHandler() func(w http.ResponseWriter, req *http.Request)
	SendDailyMessage(userID string)
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
		log.WithError(err).Fatalf("Failed to initialize chat repository")
	}
	log.Infoln("Successfully initialized chat repository")
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
		log.Errorln(err)
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
						log.Errorln(err)
						l.replyDefaultMessage(event.ReplyToken)
					}
					switch data.Intent {

					}

					rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
					defer cancel()
					if _, err = l.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).WithContext(rctx).Do(); err != nil {
						log.Errorln(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = l.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Errorln(err)
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

func (l *lineRepo) SendDailyMessage(userID string) {
	var messages []linebot.SendingMessage
	messages = make([]linebot.SendingMessage, 0, 3)
	jokeExist := true
	weatherExist := true
	newsExist := true
	userName := "User"

	jokeData, err := joke.GetRandomJoke()
	if err != nil {
		jokeExist = false
		log.WithError(err).Errorln("Failed to get joke data")
	}
	weatherData, err := weather.QueryLocation("jakarta")
	if err != nil {
		weatherExist = false
		log.WithError(err).Errorln("Failed to get weather data")
	}
	newsData, err := news.GetTopNews(news.TopNewsRequestParam{Max: 3})
	if err != nil {
		newsExist = false
		log.WithError(err).Errorln("Failed to get news data")
	}
	userData, err := l.GetUserProfile(userID)
	if err != nil {
		log.WithError(err).Errorln("Failed to get user data")
		userName = "User"
	} else {
		userName = userData.DisplayName
	}

	messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("Hello %s, here are your daily stuffs", stringformat.GetFirstWord(userName))))
	if weatherExist {
		messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("Here's your daily weather update,\n\n%s, %s:\n%s\nHumidity: %d\nTemperature: %.2f degrees C", weatherData.Name, weatherData.System.Country, weatherData.Weather[0].Description, weatherData.Details.Humidity, weatherData.Details.TemperatureCelcius)))
	}
	if jokeExist {
		messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("%s\n\n%s", jokeData.Setup, jokeData.Punchline)))
	}
	if newsExist {
		newsText := fmt.Sprintf("Top news for today:\n")
		for _, v := range newsData.Articles {
			newsText = fmt.Sprintf("%s\n%s: %s\n%s", newsText, v.Source.Name, v.Title, v.URL)
		}
		messages = append(messages, linebot.NewTextMessage(newsText))
	}

	rctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = l.Client.PushMessage(userID, messages...).WithContext(rctx).Do()
	if err != nil {
		log.WithError(err).Errorln("Failed to send message")
	}
}

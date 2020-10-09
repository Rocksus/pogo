package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/Rocksus/pogo/internal/utils/stringformat"
	"github.com/Rocksus/pogo/pkg/plugin/news"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
)

type Notifier interface {
	PushDaily(ctx context.Context) error
}

type notifier struct {
	client *linebot.Client
}

func New(client *linebot.Client) Notifier {
	return &notifier{client: client}
}

func (n *notifier) PushDaily(ctx context.Context) (err error) {
	userID := "???"
	n.sendMessage(ctx, userID)
	return
}

// TODO: refactor to use repo
func (n *notifier) sendMessage(ctx context.Context, userID string) {
	var messages []linebot.SendingMessage
	messages = make([]linebot.SendingMessage, 0, 3)
	// jokeExist := true
	// weatherExist := true
	newsExist := true
	userName := "User"

	// jokeData, err := joke.GetRandomJoke()
	// if err != nil {
	// 	jokeExist = false
	// 	log.WithError(err).Errorln("Failed to get joke data")
	// }
	// weatherData, err := weather.QueryLocation("jakarta")
	// if err != nil {
	// 	weatherExist = false
	// 	log.WithError(err).Errorln("Failed to get weather data")
	// }
	newsData, err := news.GetTopNews(news.TopNewsRequestParam{Max: 3})
	if err != nil {
		newsExist = false
		log.WithError(err).Errorln("Failed to get news data")
	}
	userData, err := n.getUserProfile(userID)
	if err != nil {
		log.WithError(err).Errorln("Failed to get user data")
		userName = "User"
	} else {
		userName = userData.DisplayName
	}

	messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("Hello %s, here are your daily stuffs", stringformat.GetFirstWord(userName))))
	// if weatherExist {
	// 	messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("Here's your daily weather update,\n\n%s, %s:\n%s\nHumidity: %d\nTemperature: %.2f degrees C", weatherData.Name, weatherData.System.Country, weatherData.Weather[0].Description, weatherData.Details.Humidity, weatherData.Details.TemperatureCelcius)))
	// }
	// if jokeExist {
	// 	messages = append(messages, linebot.NewTextMessage(fmt.Sprintf("%s\n\n%s", jokeData.Setup, jokeData.Punchline)))
	// }
	if newsExist {
		newsText := fmt.Sprintf("Top news for today:\n")
		for _, v := range newsData.Articles {
			newsText = fmt.Sprintf("%s\n%s: %s\n%s", newsText, v.Source.Name, v.Title, v.URL)
		}
		messages = append(messages, linebot.NewTextMessage(newsText))
	}

	_, err = n.client.PushMessage(userID, messages...).WithContext(ctx).Do()
	if err != nil {
		log.WithError(err).Errorln("Failed to send message")
	}
}

// TODO: move to repo
func (n *notifier) getUserProfile(userID string) (*UserData, error) {
	if userID == "" {
		return nil, fmt.Errorf("[Chat][GetUserProfile] UserID can't be empty")
	}

	rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	userData, err := n.client.GetProfile(userID).WithContext(rctx).Do()
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

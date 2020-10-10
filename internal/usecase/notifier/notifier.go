package notifier

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Rocksus/pogo/internal/repositories/pushnotif"
	"github.com/Rocksus/pogo/internal/repositories/user"
	"github.com/Rocksus/pogo/internal/utils/stringformat"
	"github.com/Rocksus/pogo/pkg/plugin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
)

type Notifier interface {
	PushDaily(ctx context.Context) error
}

type notifier struct {
	pushNotifier pushnotif.Notifier
	userRepo     user.Repository
	plugins      []plugin.DailyNotifier
}

func New(pushNotifier pushnotif.Notifier, userRepo user.Repository, plugins []plugin.DailyNotifier) Notifier {
	return &notifier{
		pushNotifier: pushNotifier,
		userRepo:     userRepo,
		plugins:      plugins,
	}
}

func (n *notifier) PushDaily(ctx context.Context) (err error) {
	userID := "???"
	n.sendMessage(ctx, userID)
	return
}

func (n *notifier) sendMessage(ctx context.Context, userID string) {
	var messages []linebot.SendingMessage
	messages = append(messages, n.createGreeting(ctx, userID))

	for _, p := range n.plugins {
		msg, err := p.GetDaily(ctx, userID)
		if err != nil {
			log.WithError(err).Errorln("failed to get daily")
			continue
		}
		messages = append(messages, msg)
	}

	err := n.pushNotifier.PushMessages(ctx, userID, messages...)
	if err != nil {
		log.WithError(err).WithField("userID", userID).Errorln("failed to push message")
	}
}

func (n *notifier) createGreeting(ctx context.Context, userID string) linebot.SendingMessage {
	helloUser, helloNoUser := randomHello()
	profile, err := n.userRepo.GetByID(ctx, userID)
	if err != nil {
		return linebot.NewTextMessage(helloNoUser)
	}

	return linebot.NewTextMessage(fmt.Sprintf(helloUser, stringformat.GetFirstWord(profile.Name)))
}

func randomHello() (withUser, noUser string) {
	choices := [][2]string{
		{
			"Hey %s, how are you doing? Here are your daily stuffs :)",
			"Hey, how are you doing? Here are your daily stuffs :)",
		},
		{
			"Hello %s, what's up? Here are some things to kick off your day",
			"Hello! Here are some things to kick off your day",
		},
		{
			"Hi %s, how's it going? Here's some stuff to start your day",
			"Hi! Here's some stuff to start your day :D",
		},
	}

	i := rand.Intn(len(choices))
	return choices[i][0], choices[i][1]
}

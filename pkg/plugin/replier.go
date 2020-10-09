package plugin

import (
	"context"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

type User struct {
	ID   string
	Name string
}

// Message follows the response structure of wit.ai.
type Message struct {
	// Sender contains the user info that sent this message. If nil, that means the sender info can't be retrieved,
	// or the sender isn't a user.
	Sender User
	// Text is the original message string
	Text string
	// Intent is the parsed message intent
	Intent string
	// Entities might contain any params required by the intent.
	Entities map[string]interface{}
	// Traits are additional message traits, such as sentiment.
	Traits map[string]interface{}
	// Timestamp is the time that the message got sent.
	Timestamp time.Time
}

type MessageReplier interface {
	Reply(ctx context.Context, message Message, replyCh chan<- linebot.SendingMessage)
}

type MessageReplierFunc func(ctx context.Context, message Message, replyCh chan<- linebot.SendingMessage)

func (f MessageReplierFunc) Reply(ctx context.Context, message Message, replyCh chan<- linebot.SendingMessage) {
	f(ctx, message, replyCh)
}

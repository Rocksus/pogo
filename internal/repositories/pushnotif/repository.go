package pushnotif

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Notifier interface {
	PushMessages(ctx context.Context, to string, messages ...linebot.SendingMessage) error
}

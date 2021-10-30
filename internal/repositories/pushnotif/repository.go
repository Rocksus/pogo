package pushnotif

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Notifier interface {
	PushMessages(ctx context.Context, to string, messages ...linebot.SendingMessage) error
}

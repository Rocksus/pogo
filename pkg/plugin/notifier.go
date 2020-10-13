package plugin

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
)

type DailyNotifier interface {
	GetDaily(ctx context.Context, recipientID string) (linebot.SendingMessage, error)
}

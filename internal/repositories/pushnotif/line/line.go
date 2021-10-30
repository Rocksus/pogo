package line

import (
	"context"

	"github.com/Rocksus/pogo/internal/repositories/pushnotif"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type pushNotifier struct {
	client *linebot.Client
}

func NewPushNotifier(client *linebot.Client) pushnotif.Notifier {
	return &pushNotifier{client: client}
}

func (p *pushNotifier) PushMessages(ctx context.Context, to string, messages ...linebot.SendingMessage) error {
	_, err := p.client.PushMessage(to, messages...).WithContext(ctx).Do()
	return err
}

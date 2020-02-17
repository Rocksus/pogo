package chat

import (
	"github.com/Rocksus/pogo/internal/repositories/interpretor"
	"github.com/line/line-bot-sdk-go/linebot"
)

type lineRepo struct {
	MasterID           string
	ChannelAccessToken string
	ChannelSecret      string
	Client             *linebot.Client
	Interpretor        interpretor.Interpretor
}

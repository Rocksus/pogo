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

type UserData struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

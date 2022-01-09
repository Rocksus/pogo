package moneysheets

import (
	"context"
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/nickylogan/go-log"

	"github.com/Rocksus/pogo/pkg/plugin"
)

type checkBalanceReplier struct {
	p *Plugin
}

func (atr checkBalanceReplier) Reply(ctx context.Context, message plugin.Message, replyCh chan<- linebot.SendingMessage) {
	var (
		responseMsg = "Sorry, I can't check that at the moment!"
		balanceType = "Unknown"
	)

	if balanceTypeRaw, ok := message.Entities["financialPlanning_balanceAccount"]; ok {
		if balanceTypeParsed, ook := balanceTypeRaw.(string); ook {
			balanceType = balanceTypeParsed
		}
	}

	amt, err := atr.p.CheckBalance(ctx, balanceType)
	if err != nil {
		log.WithError(err).Error("[CheckBalance][Reply] Failed to check balance")
	} else {
		responseMsg = fmt.Sprintf("You have %s in %s", formatAmount(amt), balanceType)
	}
	replyCh <- linebot.NewTextMessage(responseMsg)
}

func (p *Plugin) GetBalanceChecker() plugin.MessageReplier {
	return &checkBalanceReplier{
		p: p,
	}
}

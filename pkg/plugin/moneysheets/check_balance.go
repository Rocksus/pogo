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

	if balanceTypeEntityRaw, ok := message.Entities["financialPlanning_balanceAccount"]; ok {
		if balanceTypeEntity, ook := balanceTypeEntityRaw.([]interface{}); ook && len(balanceTypeEntity) > 0 {
			if balanceTypeEntityMap, oook := balanceTypeEntity[0].(map[string]interface{}); oook {
				balanceType = balanceTypeEntityMap["value"].(string)
			}
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

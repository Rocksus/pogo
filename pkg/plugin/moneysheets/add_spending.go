package moneysheets

import (
	"context"
	"math"

	"github.com/Rocksus/pogo/pkg/plugin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/nickylogan/go-log"
)

type addSpendingReplier struct {
	p *Plugin
}

func (atr addSpendingReplier) Reply(ctx context.Context, message plugin.Message, replyCh chan<- linebot.SendingMessage) {
	var (
		balanceType         = "Unknown"
		category            = ""
		responseMsg         = "Sorry, I can't process this transaction!"
		amount      float64 = 0
	)

	if balanceAccountEntityRaw, ok := message.Entities["financialPlanning_balanceAccount"]; ok {
		if balanceAccountEntity, ook := balanceAccountEntityRaw.(map[string]interface{}); ook {
			balanceType, _ = balanceAccountEntity["value"].(string)
		}
	}

	if balanceCategoryEntityRaw, ok := message.Entities["financialPlanning_category"]; ok {
		if balanceCategoryEntity, ook := balanceCategoryEntityRaw.(map[string]interface{}); ook {
			category, _ = balanceCategoryEntity["value"].(string)
		}
	}

	if amountOfMoneyEntityRaw, ok := message.Entities["amount_of_money"]; ok {
		if amountOfMoneyEntity, ook := amountOfMoneyEntityRaw.(map[string]interface{}); ook {
			// currency unit is currently not used
			amount, _ = amountOfMoneyEntity["value"].(float64)
		}
	}

	// maybe the number is classified as number
	if amount == 0 {
		if amountNumberEntityRaw, ok := message.Entities["number"]; ok {
			if amountNumberEntity, ook := amountNumberEntityRaw.(map[string]interface{}); ook {
				amount, _ = amountNumberEntity["value"].(float64)
			}
		}
	}

	if amount != 0 {
		// ensure that amount will be negative, since this is adding spending
		absAmount := math.Abs(amount)
		res, err := atr.p.AddTransaction(ctx, balanceType, category, -absAmount)
		if err != nil {
			log.WithError(err).Error("[AddTransaction][Reply] Failed to add transaction")
		} else {
			responseMsg = res
		}
	}
	replyCh <- linebot.NewTextMessage(responseMsg)
}

func (p *Plugin) GetSpendingAdder() plugin.MessageReplier {
	return &addSpendingReplier{
		p: p,
	}
}

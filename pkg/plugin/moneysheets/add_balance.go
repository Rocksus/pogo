package moneysheets

import (
	"context"
	"math"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/nickylogan/go-log"

	"github.com/Rocksus/pogo/pkg/plugin"
)

type addTransactionReplier struct {
	p *Plugin
}

func (atr addTransactionReplier) Reply(ctx context.Context, message plugin.Message, replyCh chan<- linebot.SendingMessage) {
	var (
		balanceType         = "Unknown"
		category            = ""
		responseMsg         = "Sorry, I can't process this transaction!"
		amount      float64 = 0
	)

	if balanceAccountEntityRaw, ok := message.Entities["financialPlanning_balanceAccount"]; ok {
		if balanceAccountEntity, ook := balanceAccountEntityRaw.([]interface{}); ook && len(balanceAccountEntity) > 0 {
			if balanceAccountEntityMap, oook := balanceAccountEntity[0].(map[string]interface{}); oook {
				balanceType, _ = balanceAccountEntityMap["value"].(string)
			}
		}
	}

	if balanceCategoryEntityRaw, ok := message.Entities["financialPlanning_category"]; ok {
		if balanceCategoryEntity, ook := balanceCategoryEntityRaw.([]interface{}); ook && len(balanceCategoryEntity) > 0 {
			if balanceCategoryEntityMap, oook := balanceCategoryEntity[0].(map[string]interface{}); oook {
				category, _ = balanceCategoryEntityMap["value"].(string)
			}
		}
	}

	if amountOfMoneyEntityRaw, ok := message.Entities["amount_of_money"]; ok {
		if amountOfMoneyEntity, ook := amountOfMoneyEntityRaw.([]interface{}); ook && len(amountOfMoneyEntity) > 0 {
			if amountOfMoneyEntityMap, oook := amountOfMoneyEntity[0].(map[string]interface{}); oook {
				// currency unit is currently not used
				amount, _ = amountOfMoneyEntityMap["value"].(float64)
			}
		}
	}

	// maybe the number is classified as number
	if amount == 0 {
		if amountNumberEntityRaw, ok := message.Entities["number"]; ok {
			if amountNumberEntity, ook := amountNumberEntityRaw.([]interface{}); ook && len(amountNumberEntity) > 0 {
				if amountNumberEntityMap, oook := amountNumberEntity[0].(map[string]interface{}); oook {
					amount, _ = amountNumberEntityMap["value"].(float64)
				}
			}
		}
	}

	if amount != 0 {
		// force positive amount
		absAmount := math.Abs(amount)
		res, err := atr.p.AddTransaction(ctx, balanceType, category, absAmount)
		if err != nil {
			log.WithError(err).Error("[AddTransaction][Reply] Failed to add transaction")
		} else {
			responseMsg = res
		}
	}
	replyCh <- linebot.NewTextMessage(responseMsg)
}

func (p *Plugin) GetTransactionAdder() plugin.MessageReplier {
	return &addTransactionReplier{
		p: p,
	}
}

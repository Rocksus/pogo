package moneysheets

import (
	"context"
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/nickylogan/go-log"
	"google.golang.org/api/sheets/v4"

	"github.com/Rocksus/pogo/pkg/plugin"
)

type addTransactionReplier struct {
	p *Plugin
}

func (atr addTransactionReplier) Reply(ctx context.Context, message plugin.Message, replyCh chan<- linebot.SendingMessage) {
	var (
		balanceType       = "Unknown"
		category          = ""
		responseMsg       = "Sorry, I can't process this transaction!"
		amount      int64 = 0
		err         error
	)
	if balanceTypeRaw, ok := message.Entities["financialPlanning_balanceAccount"]; ok {
		if balanceTypeParsed, ook := balanceTypeRaw.(string); ook {
			balanceType = balanceTypeParsed
		}
	}

	if balanceCategoryRaw, ok := message.Entities["financialPlanning_category"]; ok {
		if balanceCategoryParsed, ook := balanceCategoryRaw.(string); ook {
			category = balanceCategoryParsed
		}
	}

	if amountRaw, ok := message.Entities["wit/amount_of_money"]; ok {
		if amountParsed, ook := amountRaw.(string); ook {
			amount, err = processRawAmount(amountParsed)
			if err != nil {
				log.WithError(err).Error("[AddTransaction][Reply] Failed to process amount")
			}
		}
	}

	if amount != 0 {
		res, err := atr.p.AddTransaction(ctx, balanceType, category, amount)
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

func (p *Plugin) AddTransaction(ctx context.Context, balanceType, category string, amount int64) (resp string, err error) {
	_, err = p.sheet.Spreadsheets.Values.Append(
		p.sheetID,
		sheetRange,
		&sheets.ValueRange{
			MajorDimension: "ROWS",
			Values: [][]interface{}{
				{time.Now().Format(dateFormat), amount, nil, balanceType, category},
			},
		},
	).
		ValueInputOption("USER_ENTERED").
		Context(ctx).Do()
	if err != nil {
		return "Failed to insert transaction", err
	}
	return fmt.Sprintf("Added new transaction of %d for %s", amount, balanceType), nil
}

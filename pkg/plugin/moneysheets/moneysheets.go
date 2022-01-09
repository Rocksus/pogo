package moneysheets

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/Rocksus/pogo/configs"
)

type Plugin struct {
	sheet   *sheets.Service
	sheetID string
}

func GetScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/spreadsheets",
	}
}

func InitPlugin(cfg configs.MoneySheetsConfig, gclient *http.Client) (*Plugin, error) {
	sheet, err := sheets.NewService(context.Background(), option.WithHTTPClient(gclient))
	if err != nil {
		return nil, err
	}

	return &Plugin{
		sheet:   sheet,
		sheetID: cfg.SheetID,
	}, nil
}

func (p *Plugin) AddTransaction(ctx context.Context, balanceType, category string, amount float64) (resp string, err error) {
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
	return fmt.Sprintf("Added new transaction of %f for %s", amount, balanceType), nil
}

func (p *Plugin) CheckBalance(ctx context.Context, balanceType string) (float64, error) {
	var amount float64
	// balanceRange only applies if you are using the provided google sheets template
	balanceRange := "Bank Management!J1:2"
	res, err := p.sheet.Spreadsheets.Values.Get(p.sheetID, balanceRange).ValueRenderOption("UNFORMATTED_VALUE").Context(ctx).Do()
	if err != nil {
		return 0, err
	}

	if len(res.Values) == 0 {
		return 0, nil
	}

	loweredBalanceType := strings.ToLower(balanceType)

	for i := range res.Values[0] {
		if sheetBalanceType, ok := res.Values[0][i].(string); ok {
			if strings.ToLower(sheetBalanceType) == loweredBalanceType {
				amount = res.Values[1][i].(float64)
			}
		}
	}

	return amount, nil
}

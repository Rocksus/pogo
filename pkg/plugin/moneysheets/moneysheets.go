package moneysheets

import (
	"context"
	"net/http"

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

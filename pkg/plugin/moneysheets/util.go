package moneysheets

import "strconv"

func formatAmount(amt float64) string {
	// TODO: adjust currency
	return "Rp." + strconv.FormatFloat(amt, 'f', 2, 64)
}

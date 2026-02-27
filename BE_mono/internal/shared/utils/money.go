package utils

import "fmt"

func ToAmountString(cents int64) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}

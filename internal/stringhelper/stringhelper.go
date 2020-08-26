package stringhelper

import (
	"fmt"
	"strconv"
)

func ShortScale(digit int64) string {
	if digit >= 1000000000 {
		return fmt.Sprintf("%.1fb", float64(digit)/1000000000)
	}
	if digit >= 1000000 {
		return fmt.Sprintf("%.1fm", float64(digit)/1000000)
	}
	if digit >= 1000 {
		return fmt.Sprintf("%.1fk", float64(digit)/1000)
	}
	return strconv.FormatInt(digit, 10)
}

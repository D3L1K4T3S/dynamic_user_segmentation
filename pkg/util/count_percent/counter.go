package count_percent

import "math"

// CheckCount метод, который по выдает true или false
// в зависимости от того, входит ли номер пользователя в данный процент
func CheckCount(count int, percent float64) bool {
	number := int(math.Round(1 / percent * 100))
	if count%number == 0 {
		return true
	}
	return false
}

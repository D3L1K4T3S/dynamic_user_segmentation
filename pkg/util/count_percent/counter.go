package count_percent

import "math"

func CheckCount(count int, percent float64) bool {
	number := int(math.Round(1 / percent * 100))
	if count%number == 0 {
		return true
	}
	return false
}

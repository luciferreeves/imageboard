package format

import "fmt"

func Count(count int64) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	} else if count < 1000000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	} else {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	}
}

func Int64ToString(value int64) string {
	if value < 0 {
		return fmt.Sprintf("-%d", -value)
	}
	return fmt.Sprintf("%d", value)
}

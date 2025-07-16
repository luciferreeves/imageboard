package format

import "fmt"

func FileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	val := float64(size) / float64(div)
	unitStr := "KMGTPE"[exp : exp+1]
	if val == float64(int64(val)) {
		return fmt.Sprintf("%d %sB", int64(val), unitStr)
	}
	return fmt.Sprintf("%.2f %sB", val, unitStr)
}

func Count(count int64) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	} else if count < 1000000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	} else {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	}
}

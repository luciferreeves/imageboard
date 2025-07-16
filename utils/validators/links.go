package validators

import "strings"

func GetRefererForURL(url string) string {
	switch {
	case strings.Contains(url, "i.pximg.net") || strings.Contains(url, "pixiv.net"):
		return "https://www.pixiv.net/"
	default:
		return ""
	}
}

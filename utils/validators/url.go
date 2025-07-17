package validators

import (
	"regexp"
	"strings"
)

func IsValidURL(url string) bool {
	if url == "" {
		return false
	}
	if len(url) > 2048 { // Common max URL length
		return false
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	pattern := `^(http|https)://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=]+$`
	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		return false
	}

	return matched
}

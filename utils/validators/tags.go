package validators

import "regexp"

func IsValidTagName(tag string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, tag)
	return match
}

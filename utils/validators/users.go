package validators

import (
	"regexp"
	"slices"
	"strings"
)

func IsValidUsername(username string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return match
}

func IsReservedUsername(username string) bool {
	reserved := []string{
		"admin", "administrator", "mod", "moderator", "janitor",
		"api", "www", "mail", "email", "support", "help",
		"about", "contact", "privacy", "terms", "tos",
		"null", "undefined", "system", "bot", "guest",
		"login", "register", "signup", "signin", "logout",
		"profile", "settings", "shifoo", "deleted",
	}

	lowerUsername := strings.ToLower(username)
	return slices.Contains(reserved, lowerUsername)
}

func IsValidEmail(email string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return match
}

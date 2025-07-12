package processors

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SitePreferences struct {
	SidebarWidth     string `json:"sidebar_width"`
	MainContentWidth string `json:"main_content_width"`
	H1FontSize       string `json:"h1_font_size"`
	BodyFontSize     string `json:"body_font_size"`
	SmallFontSize    string `json:"small_font_size"`
	PostsPerPage     int    `json:"posts_per_page"`
}

func PreferencesContextProcessor(context *fiber.Ctx) error {
	defaultPreferences := SitePreferences{
		SidebarWidth:     "180px",
		MainContentWidth: "1200px",
		H1FontSize:       "16px",
		BodyFontSize:     "13px",
		SmallFontSize:    "11px",
		PostsPerPage:     42,
	}

	preferences := defaultPreferences

	preferencesCookie := context.Cookies("preferences")
	if preferencesCookie != "" {
		_ = json.Unmarshal([]byte(preferencesCookie), &preferences)
	}

	bytes, err := json.Marshal(preferences)
	if err == nil {
		context.Cookie(&fiber.Cookie{
			Name:     "preferences",
			Value:    string(bytes),
			Path:     "/",
			SameSite: fiber.CookieSameSiteLaxMode,
		})
	}

	context.Locals("Preferences", preferences)
	context.Locals("PreferencesCSS", preferencesToCSS(preferences))
	return context.Next()
}

func preferencesToCSS(preferences SitePreferences) string {
	return fmt.Sprintf(`
	<style>
		main {
			width: %s;
		}
		body {
			font-size: %s;
		}
		h1 {
			font-size: %s;
		}
		small {
			font-size: %s;
		}
		.sidebar {
			width: %s;
		}
	</style>`,
		preferences.MainContentWidth,
		preferences.BodyFontSize,
		preferences.H1FontSize,
		preferences.SmallFontSize,
		preferences.SidebarWidth,
	)
}

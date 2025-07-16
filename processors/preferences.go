package processors

import (
	"encoding/json"
	"fmt"
	"imageboard/config"

	"github.com/gofiber/fiber/v2"
)

func PreferencesContextProcessor(context *fiber.Ctx) error {
	defaultPreferences := config.SitePreferences{
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

func preferencesToCSS(preferences config.SitePreferences) string {
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

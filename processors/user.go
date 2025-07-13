package processors

import (
	"imageboard/database"
	"imageboard/models"
	"imageboard/session"

	"github.com/gofiber/fiber/v2"
)

func UserContextProcessor(ctx *fiber.Ctx) error {
	var user *models.User

	sess, err := session.Store.Get(ctx)
	if err == nil {
		var userID uint
		if id := sess.Get("user_id"); id != nil {
			switch v := id.(type) {
			case uint:
				userID = v
			case int:
				userID = uint(v)
			case int64:
				userID = uint(v)
			case float64:
				userID = uint(v)
			}
		}

		if userID != 0 {
			dbUser, err := database.GetUserByID(userID)
			if err == nil && dbUser != nil {
				dbUser.UpdateLastUserActivity(database.DB)
				user = dbUser
			}
		}
	}

	ctx.Locals("User", user)
	ctx.Locals("IsAuthenticated", user != nil)

	return ctx.Next()
}

package session

import (
	"imageboard/config"
	"log"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v2"
)

var Store *session.Store

func init() {
	storage := postgres.New(postgres.Config{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Username: config.Database.Username,
		Password: config.Database.Password,
		Database: config.Database.DatabaseName,
		Table:    "sessions",
		SSLMode:  config.Database.SSLMode,
	})

	Store = session.New(session.Config{
		Storage:        storage,
		Expiration:     config.Session.Expiration,
		KeyLookup:      "cookie:" + config.Session.CookieName,
		CookieName:     config.Session.CookieName,
		CookieDomain:   config.Session.CookieDomain,
		CookiePath:     config.Session.CookiePath,
		CookieSecure:   config.Session.CookieSecure,
		CookieSameSite: config.Session.CookieSameSite,
		CookieHTTPOnly: true,
	})

	log.Println("Session store initialized successfully")
}

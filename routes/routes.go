package routes

import (
	"bloom/config"
	"bloom/routes/observations"
	"bloom/routes/taxon"
	"bloom/routes/users"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/contrib/jwt"
)

func Init(app *fiber.App) {
	usersG := app.Group("/users")
	taxonG := app.Group("/taxon")
	observationsG := app.Group("/observationsG")

	usersG.Post("/user", users.Insert)
	usersG.Post("/login", users.Login)

	taxonG.Get("/search", taxon.Search)

	observationsG.Get("/observation/specific", observations.Get)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.JWTSigningKey)},
	}))

	usersG.Delete("/user", users.Delete)
	usersG.Patch("/user/location", users.UpdateLocation)
	usersG.Patch("/user/password", users.UpdatePassword)

	observationsG.Get("/observation", observations.GetSelf)
	observationsG.Post("/observation", observations.Insert)
	observationsG.Patch("/observation", observations.Patch)
	observationsG.Delete("/observation", observations.Delete)
}

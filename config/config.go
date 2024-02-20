package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const ListenAddress = ":3000"

var FiberConfig = fiber.Config{
	Prefork:       true,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "Fiber",
	AppName:       "Bloom",
}

const JWTSigningKey = "ChangeMe!!!!!"

const MongoDBURI = "mongodb://localhost:27017"
const MongoDBConnTimeout = 20 * time.Second
const MongoDBName = "bloom"

const ReturnFirstXResults = 50

const PasswordHashSalt = "CHANGEME!"

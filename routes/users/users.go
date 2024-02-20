package users

import (
	"bloom/config"
	"bloom/db"
	"bloom/db/models/users"
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type insertUserModel struct {
	Username    string
	Email       string
	Password    string
	Lat         float64
	Long        float64
	Permissions int
}

func Insert(c *fiber.Ctx) error {
	insertUser := new(insertUserModel)
	dec := json.NewDecoder(bytes.NewReader(c.Body()))
	dec.DisallowUnknownFields()
	err := dec.Decode(insertUser)
	if err != nil {
		c.JSON(fiber.Map{"error": "error with JSON parsing"})
		return c.SendStatus(400)
	}
	usersColl := db.DB().Collection("users")
	if insertUser.Username == "" {
		c.JSON(fiber.Map{"error": "username cannot be empty"})
		return c.SendStatus(400)
	}
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(insertUser.Password+config.PasswordHashSalt), 14)
	if err != nil {
		c.JSON(fiber.Map{"error": "invalid password. is it longer than 72 bytes?"})
		return c.SendStatus(400) // Usually because the password is too long
	}
	existingUser := usersColl.FindOne(context.Background(), bson.M{"username": insertUser.Username})
	if existingUser.Err() == nil {
		c.JSON(fiber.Map{"error": "user already exists"})
		return c.SendStatus(400) // Usually because the password is too long

	}

	if insertUser.Permissions != 4 {
		c.JSON(fiber.Map{"error": "user permissions invalid"})
		return c.SendStatus(401)
	}

	entry := users.Entry{
		Username:    insertUser.Username,
		Email:       insertUser.Email,
		Permissions: insertUser.Permissions,
		Password:    string(passwordBytes),
		DateCreated: time.Now().Unix(),
		Lat:         insertUser.Lat,
		Long:        insertUser.Long,
		Garden:      []string{},
	}
	_, err = usersColl.InsertOne(context.Background(), entry)
	if err != nil {
		c.JSON(fiber.Map{"error": "error with inserting into database"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null"})
	return c.SendStatus(200)
}

func Login(c *fiber.Ctx) error {
	usersColl := db.DB().Collection("users")

	user := c.FormValue("username")
	pass := c.FormValue("password")

	// Throws Unauthorized error
	existingUser := usersColl.FindOne(context.Background(), bson.M{"username": user})
	if existingUser.Err() != nil {
		c.JSON(fiber.Map{"error": "user doesn't exist"})
		return c.SendStatus(fiber.StatusUnauthorized)

	}
	userO := new(users.Entry)
	err := existingUser.Decode(userO)

	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	bcrypt.CompareHashAndPassword([]byte(userO.Password), []byte(pass))

	// Create the Claims
	claims := jwt.MapClaims{
		"username": userO.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.JWTSigningKey))

	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"error": "null", "token": t})
}

func Delete(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)

	userID, err := GetUserIDFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	usersColl := db.DB().Collection("users")
	_, err = usersColl.DeleteOne(context.Background(), bson.M{"_id": userID})
	if err != nil {
		c.JSON(fiber.Map{"error": "error with deleting user"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null"})
	return c.SendStatus(200)
}

func UpdatePassword(c *fiber.Ctx) error {
	password := c.Query("password")

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)

	userID, err := GetUserIDFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password+config.PasswordHashSalt), 14)
	if err != nil {
		c.JSON(fiber.Map{"error": "invalid password. is it longer than 72 bytes?"})
		return c.SendStatus(400) // Usually because the password is too long
	}

	usersColl := db.DB().Collection("users")
	entry := users.Entry{
		Password: string(passwordBytes),
	}
	_, err = usersColl.UpdateByID(context.Background(), userID, entry)
	if err != nil {
		c.JSON(fiber.Map{"error": "error with updating entry in database"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null"})
	return c.SendStatus(200)
}

func UpdateLocation(c *fiber.Ctx) error {
	latS := c.Query("lat")
	longS := c.Query("long")
	lat, err := strconv.ParseFloat(latS, 64)
	if err != nil {
		c.JSON(fiber.Map{"error": "error in parsing float64"})
		return c.SendStatus(400)
	}
	long, err := strconv.ParseFloat(longS, 64)
	if err != nil {
		c.JSON(fiber.Map{"error": "error in parsing float64"})
		return c.SendStatus(400)
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)

	userID, err := GetUserIDFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	usersColl := db.DB().Collection("users")
	entry := users.Entry{
		Lat:  lat,
		Long: long,
	}
	_, err = usersColl.UpdateByID(context.Background(), userID, entry)
	if err != nil {
		c.JSON(fiber.Map{"error": "error with updating entry in database"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null"})
	return c.SendStatus(200)
}

type IDReceiv struct {
	ObjectID primitive.ObjectID `json:"-" bson:"_id"`
}

func GetUserIDFromUsername(username string) (string, error) {
	usersColl := db.DB().Collection("users")
	existingUser := usersColl.FindOne(context.Background(), bson.M{"username": username})
	user := new(IDReceiv)
	err := existingUser.Decode(user)
	return user.ObjectID.String(), err
}

func GetUserFromUsername(username string) (*users.Entry, error) {
	usersColl := db.DB().Collection("users")
	existingUser := usersColl.FindOne(context.Background(), bson.M{"username": username})
	user := new(users.Entry)
	err := existingUser.Decode(user)
	return user, err
}

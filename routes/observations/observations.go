package observations

import (
	"bloom/db"
	"bloom/db/models/observations"
	"bloom/routes/users"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type insertObservationModel struct {
	TaxonID        string
	Lat            float64
	Long           float64
	BloomStartDate int64
	BloomPeakDate  int64
	BloomEndDate   int64
}

type patchObservationModel struct {
	ID             string
	TaxonID        string
	Lat            float64
	Long           float64
	BloomStartDate int64
	BloomPeakDate  int64
	BloomEndDate   int64
}

func Insert(c *fiber.Ctx) error {
	insertObservation := new(insertObservationModel)
	dec := json.NewDecoder(bytes.NewReader(c.Body()))
	dec.DisallowUnknownFields()
	err := dec.Decode(insertObservation)
	if err != nil {
		c.JSON(fiber.Map{"error": "error with JSON parsing"})
		return c.SendStatus(400)
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)

	userID, err := users.GetUserIDFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	if math.Abs(insertObservation.Lat) > 90 || math.Abs(insertObservation.Long) > 180 {
		c.JSON(fiber.Map{"error": "invalid coordinates"})
		return c.SendStatus(400)
	}

	observationsColl := db.DB().Collection("observations")

	entry := observations.Entry{
		UserID:         userID,
		TaxonID:        insertObservation.TaxonID,
		Lat:            insertObservation.Lat,
		Long:           insertObservation.Long,
		BloomStartDate: insertObservation.BloomStartDate,
		BloomPeakDate:  insertObservation.BloomPeakDate,
		BloomEndDate:   insertObservation.BloomEndDate,
		DateModified:   time.Now().Unix(),
		DateCreated:    time.Now().Unix(),
	}
	_, err = observationsColl.InsertOne(context.Background(), entry)
	if err != nil {
		c.JSON(fiber.Map{"error": "error with inserting into database"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null"})
	return c.SendStatus(200)
}

func Patch(c *fiber.Ctx) error {
	patchObservation := new(patchObservationModel)
	dec := json.NewDecoder(bytes.NewReader(c.Body()))
	dec.DisallowUnknownFields()
	err := dec.Decode(patchObservation)
	if err != nil {
		c.JSON(fiber.Map{"error": "error with JSON parsing"})
		return c.SendStatus(400)
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)

	userID, err := users.GetUserIDFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})

		return c.SendStatus(501)
	}

	userData, err := users.GetUserFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	observationsColl := db.DB().Collection("observations")

	objectId := new(primitive.ObjectID)
	err = objectId.UnmarshalText([]byte(patchObservation.ID))
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}
	existingObs := observationsColl.FindOne(context.Background(), bson.M{"_id": objectId})
	observation := new(observations.Entry)
	err = existingObs.Decode(observation)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	if userData.Permissions != 0 {
		if observation.UserID != userID {
			c.JSON(fiber.Map{"error": "you tried to change an observation that isn't yours"})
			return c.SendStatus(401)
		}
	}
	if math.Abs(patchObservation.Lat) > 90 || math.Abs(patchObservation.Long) > 180 {
		c.JSON(fiber.Map{"error": "invalid coordinates"})
		return c.SendStatus(400)
	}
	entry := observations.Entry{
		UserID:         userID,
		TaxonID:        patchObservation.TaxonID,
		Lat:            patchObservation.Lat,
		Long:           patchObservation.Long,
		BloomStartDate: patchObservation.BloomStartDate,
		BloomPeakDate:  patchObservation.BloomPeakDate,
		BloomEndDate:   patchObservation.BloomEndDate,
		DateModified:   time.Now().Unix(),
		DateCreated:    observation.DateCreated,
	}
	_, err = observationsColl.ReplaceOne(context.Background(), bson.M{"_id": objectId}, entry)
	if err != nil {
		log.Error(err)
		c.JSON(fiber.Map{"error": "error with updating entry in database"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null"})
	return c.SendStatus(200)
}

func Delete(c *fiber.Ctx) error {
	delID := c.Query("id")

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)

	userID, err := users.GetUserIDFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	userData, err := users.GetUserFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	observationsColl := db.DB().Collection("observations")
	objectId := new(primitive.ObjectID)
	err = objectId.UnmarshalText([]byte(delID))
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}
	existingObs := observationsColl.FindOne(context.Background(), bson.M{"_id": objectId})
	observation := new(observations.Entry)
	err = existingObs.Decode(observation)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	if userData.Permissions != 0 {
		if observation.UserID != userID {
			c.JSON(fiber.Map{"error": "you tried to delete an observation that isn't yours"})
			return c.SendStatus(401)
		}
	}

	res, err := observationsColl.DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		c.JSON(fiber.Map{"error": "error with deleting entry in database"})
		return c.SendStatus(501)
	}
	if res.DeletedCount == 0 {
		c.JSON(fiber.Map{"error": "didn't find an element to delete"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "no error."})
	return c.SendStatus(200)
}

func GetSelf(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)

	userID, err := users.GetUserIDFromUsername(name)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}
	observationsColl := db.DB().Collection("observations")

	opts := options.Find().SetSort(bson.M{"dateModified": 1})
	searchResultsCur, err := observationsColl.Find(context.Background(), bson.M{"userid": userID}, opts)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		fmt.Println("here")
		return c.SendStatus(501)
	}

	var searchResults []observations.EntryIDED
	if err = searchResultsCur.All(context.TODO(), &searchResults); err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null", "data": searchResults})
	return c.SendStatus(200)
}

func Get(c *fiber.Ctx) error {
	taxonID := c.Query("id")

	observationsColl := db.DB().Collection("observations")

	opts := options.Find().SetSort(bson.M{"dateModified": 1})
	searchResultsCur, err := observationsColl.Find(context.Background(), bson.M{"taxonid": taxonID}, opts)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		fmt.Println("here")
		return c.SendStatus(501)
	}

	var searchResults []observations.Entry
	if err = searchResultsCur.All(context.TODO(), &searchResults); err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null", "data": searchResults})
	return c.SendStatus(200)
}

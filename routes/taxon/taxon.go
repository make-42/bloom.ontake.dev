package taxon

import (
	"bloom/config"
	"bloom/db"
	"bloom/db/models/taxon"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Search(c *fiber.Ctx) error {
	searchName := c.Query("q")

	taxonColl := db.DB().Collection("taxon")
	opts := options.Find().SetSort(bson.M{"scientificname": 1})
	searchResultsCur, err := taxonColl.Find(context.Background(), bson.M{"scientificname": bson.M{"$regex": "^(?i)" + searchName + "*"}, "taxonomicstatus": bson.M{"$ne": "Unchecked"}}, opts)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	var searchResults []taxon.EntryIDED
	if err = searchResultsCur.All(context.TODO(), &searchResults); err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null", "data": searchResults[:min(config.ReturnFirstXResults, len(searchResults))]})
	return c.SendStatus(200)
}

func ByID(c *fiber.Ctx) error {
	id := c.Query("id")

	taxonColl := db.DB().Collection("taxon")

	objectId := new(primitive.ObjectID)
	err := objectId.UnmarshalText([]byte(id))
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}
	searchResult := taxonColl.FindOne(context.Background(), bson.M{"_id": objectId})

	res := new(taxon.Entry)

	searchResult.Decode(res)

	c.JSON(fiber.Map{"error": "null", "data": res})
	return c.SendStatus(200)
}

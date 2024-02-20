package taxon

import (
	"bloom/db"
	"bloom/db/models/taxon"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Search(c *fiber.Ctx) error {
	searchName := c.Query("q")

	taxonColl := db.DB().Collection("taxon")
	opts := options.Find().SetSort(bson.M{"scientificname": 1})
	searchResultsCur, err := taxonColl.Find(context.Background(), bson.M{"scientificname": bson.M{"$regex": "^" + searchName + "*"}, "taxonomicstatus": bson.M{"$ne": "Unchecked"}}, opts)
	if err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		fmt.Println("here")
		return c.SendStatus(501)
	}

	var searchResults []taxon.Entry
	if err = searchResultsCur.All(context.TODO(), &searchResults); err != nil {
		c.JSON(fiber.Map{"error": "internal server error"})
		return c.SendStatus(501)
	}

	c.JSON(fiber.Map{"error": "null", "data": searchResults})
	return c.SendStatus(200)
}

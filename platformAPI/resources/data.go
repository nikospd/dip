package resources

import (
	"context"
	"dev.com/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func GetStorageData(c echo.Context, client *mongo.Client, resourcesDb string, dataDb string, storageCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")
	//Get storage from id
	//Check if storage belongs to this userId or if it is shared to it
	findQuery := bson.D{
		{"_id", storageId},
		{"$or", bson.A{
			bson.D{{"user_id", userId}},
			bson.D{{"shared_with_id", bson.D{{"$in", bson.A{userId}}}}}},
		}}
	storageCollection := client.Database(resourcesDb).Collection(storageCol)
	one := storageCollection.FindOne(context.TODO(), findQuery)
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad gateway"})
	}
	var storage utils.Storage
	err := one.Decode(&storage)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get storage"})
	}
	//Get data from that storage
	dataCollection := client.Database(dataDb).Collection(storageId)
	opts := options.Find()
	opts.Projection = bson.D{{"_id", 0}}
	cur, dferr := dataCollection.Find(context.TODO(), bson.D{}, opts)
	if dferr != nil {
		if dferr == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var dataTable []map[string]interface{}
	cur.All(context.TODO(), &dataTable)
	return c.JSON(http.StatusOK, dataTable)
}

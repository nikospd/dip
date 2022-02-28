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

/*
	This module is about the filters a user may apply at incoming data for
	selectively process and store.
*/

func GetStorageFilter(c echo.Context, client *mongo.Client, dataDb string, resourcesDb string, filterCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")

	//Get first document from that storage
	collection := client.Database(dataDb).Collection(storageId)
	opts := options.FindOne()
	opts.Projection = bson.D{{"_id", 0}}
	opts.SetSort(bson.D{{"arrived_at", 1}})
	cur := collection.FindOne(context.TODO(), bson.D{{"user_id", userId}}, opts)
	if cur.Err() != nil {
		if cur.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "No documents at this storage yet"})
		}
	}
	var document map[string]interface{}
	err := cur.Decode(&document)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get document"})
	}

	//Get the existing filter for this storage
	collection = client.Database(resourcesDb).Collection(filterCol)
	cur = collection.FindOne(context.TODO(), bson.D{{"user_id", userId}, {"storage_id", storageId}})
	if cur.Err() != nil {
		if cur.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusOK, echo.Map{"msg": "No filter found", "document": document, "filter": bson.A{}})
		}
	}
	var filter utils.StorageFilter
	err = cur.Decode(&filter)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get filter"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK", "document": document, "filter": filter})
}

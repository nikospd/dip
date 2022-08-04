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
	"time"
)

/*
	This module is about the filters a user may apply at incoming data for
	selectively process and store.
*/

func CreateStorageFilter(c echo.Context, client *mongo.Client, db string, filterCol string) error {
	userId, storageId := utils.GetRequestIds(c)
	filter := new(utils.StorageFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get attributes"})
	}
	if len(filter.Attributes) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Empty filter attributes is not accepted"})
	}
	//Sanitize data
	token := utils.CreateRandomHash(20)
	filter.FilterId = token
	filter.UserId = userId
	filter.StorageId = storageId
	filter.CreatedAt = time.Now()
	filter.ModifiedAt = time.Time{}
	filter.Print()
	//Insert document
	collection := client.Database(db).Collection(filterCol)
	_, err := collection.InsertOne(context.TODO(), filter)
	if mongo.IsDuplicateKeyError(err) {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Filter already exist for this storage"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create resource"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "Storage filter created", "id": token})
}

func GetStorageFilter(c echo.Context, client *mongo.Client, dataDb string, resourcesDb string, filterCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")

	//Get first document from that storage
	collection := client.Database(dataDb).Collection(storageId)
	opts := options.FindOne()
	opts.Projection = bson.D{{"_id", 0}}
	opts.SetSort(bson.D{{"arrived_at", -1}})
	cur := collection.FindOne(context.TODO(), bson.D{{"user_id", userId}}, opts)
	var document map[string]interface{}
	if cur.Err() == nil {
		err := cur.Decode(&document)
		if err != nil {
			return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get document"})
		}
	}

	//Get the existing filter for this storage or integration
	integrationId := c.QueryParam("integrationId")
	collection = client.Database(resourcesDb).Collection(filterCol)
	if len(integrationId) > 0 {
		cur = collection.FindOne(context.TODO(), bson.D{{"user_id", userId}, {"storage_id", integrationId}})
	} else {
		cur = collection.FindOne(context.TODO(), bson.D{{"user_id", userId}, {"storage_id", storageId}})
	}
	if cur.Err() != nil {
		if cur.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusPartialContent, echo.Map{"msg": "No filter found", "document": document, "filter": bson.A{}})
		}
	}
	var filter utils.StorageFilter
	err := cur.Decode(&filter)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get filter"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK", "document": document, "filter": filter})
}

func UpdateStorageFilter(c echo.Context, client *mongo.Client, db string, filterCol string) error {
	userId, filterId := utils.GetRequestIds(c)
	filter := new(utils.StorageFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get attributes"})
	}
	if len(filter.Attributes) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Empty filter attributes is not accepted"})
	}
	//Sanitize data
	filter.FilterId = ""
	filter.UserId = ""
	filter.StorageId = ""
	filter.CreatedAt = time.Time{}
	filter.ModifiedAt = time.Now()
	//Update the document
	collection := client.Database(db).Collection(filterCol)
	one, err := collection.UpdateOne(context.TODO(), bson.D{{"user_id", userId}, {"_id", filterId}},
		bson.D{{"$set", filter}})
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Filter not found"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to update filter"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}
func DeleteStorageFilter(c echo.Context, client *mongo.Client, db string, filterCol string) error {
	userId, filterId := utils.GetRequestIds(c)
	collection := client.Database(db).Collection(filterCol)
	one, err := collection.DeleteOne(context.TODO(), bson.D{{"_id", filterId}, {"user_id", userId}})
	if one.DeletedCount == 0 {
		return c.JSON(http.StatusConflict, echo.Map{"msg": "No filter deleted"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to delete filter"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

package resources

import (
	"context"
	"dev.com/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

/*
	This module made for the management of the properties of a user.
	Like devices, applications and data storages (buckets)
	If it gets huge, I will separate each into different files
*/

func CreateApplication(c echo.Context, client *mongo.Client, db string, storageCol string, appCol string) error {
	//Request body has "description & persistRaw"
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	app := new(utils.Application)
	if err := c.Bind(app); err != nil {
		return err
	}
	app.AppId = utils.CreateRandomHash(20)
	if app.PersistRaw {
		if app.RawStorageId == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"msg": "RawStorageId not provided but raw persistence was selected"})
		}
		storageCollection := client.Database(db).Collection(storageCol)
		filter := bson.D{{"_id", app.RawStorageId}, {"user_id", userId}}
		updateQuery := bson.D{{"$set", bson.D{{"app_id", app.AppId}}}}
		oneStorage, err := storageCollection.UpdateOne(context.TODO(), filter, updateQuery)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
		}
		if oneStorage.MatchedCount == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Storage Id that provided, does not match"})
		}
	}
	app.UserId = userId
	app.CreatedAt = time.Now()
	collection := client.Database(db).Collection(appCol)
	_, err := collection.InsertOne(context.TODO(), app)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create resource"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "Application created"})
}

func UpdateApplication(c echo.Context, client *mongo.Client, db string, appCol string) error {
	//Request body has "description & persistRaw"
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	appId := c.Param("id")
	app := new(utils.Application)
	if err := c.Bind(app); err != nil {
		return err
	}
	collection := client.Database(db).Collection(appCol)
	filter := bson.D{{"_id", appId}, {"user_id", userId}}
	//Sanitize data
	app.AppId = ""
	app.UserId = ""
	app.CreatedAt = time.Time{}
	app.ModifiedAt = time.Now()
	//Create update query
	var updateFields bson.D
	tmpFields, _ := bson.Marshal(app)
	unmarshalErr := bson.Unmarshal(tmpFields, &updateFields)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	updateQuery := bson.D{{"$set", updateFields}}
	//Update the database
	one, err := collection.UpdateOne(context.TODO(), filter, updateQuery)
	//Handle errors and respond
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Application not found"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Application not modified"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func GetApplicationsByUser(c echo.Context, client *mongo.Client, db string, appCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	collection := client.Database(db).Collection(appCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"user_id", userId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var appTable []utils.Application
	err = cur.All(context.TODO(), &appTable)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Error at getting the applications"})
	}
	return c.JSON(http.StatusOK, appTable)
}

func DeleteApplication(c echo.Context, client *mongo.Client, db string, appCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	appId := c.Param("id")
	collection := client.Database(db).Collection(appCol)
	one, err := collection.DeleteOne(context.TODO(), bson.D{
		{"_id", appId},
		{"user_id", userId}})
	if one.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Token not deleted"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	//TODO: Delete the storages that are assigned in this app as well
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

//TODO:

func AddStorageToApplication(c echo.Context, client *mongo.Client) error {
	return nil
}

//TODO:

func RemoveStorageFromApplication(c echo.Context, client *mongo.Client) error {
	return nil
}

func CreateStorage(c echo.Context, client *mongo.Client, db string, storageCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storage := new(utils.Storage)
	if err := c.Bind(storage); err != nil {
		return err
	}
	switch storage.Type {
	case "cloudMongo":
		fmt.Println("Case simple mongo. No further actions required")
	case "proprietaryMongo":
		return c.JSON(http.StatusNotImplemented, echo.Map{"msg": "Storage option not implemented yet."})
	default:
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Wrong type of storage"})
	}
	//Init data
	storage.StorageId = utils.CreateRandomHash(20)
	storage.UserId = userId
	storage.CreatedAt = time.Now()
	storage.SharedWithId = []string{}
	collection := client.Database(db).Collection(storageCol)
	_, err := collection.InsertOne(context.TODO(), storage)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create resource"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "Storage created", "id": storage.StorageId})
}

//TODO: get, update, delete, share, stopSharing.

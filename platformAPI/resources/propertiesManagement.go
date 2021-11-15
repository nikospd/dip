package resources

import (
	"context"
	"dev.com/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	switch app.SourceType {
	case "push":
		fmt.Println("New push mechanism")
	case "pull":
		fmt.Println("New pull mechanism")
	default:
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Wrong type source"})
	}
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
	} else {
		app.RawStorageId = ""
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

func GetApplicationById(c echo.Context, client *mongo.Client, db string, appCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	appId := c.Param("id")
	collection := client.Database(db).Collection(appCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"_id", appId},
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
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Application not deleted"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	//TODO: Delete the storages that are assigned in this app as well
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
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

func GetStorage(c echo.Context, client *mongo.Client, db string, storageCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")

	collection := client.Database(db).Collection(storageCol)
	one := collection.FindOne(context.TODO(), bson.D{
		{"user_id", userId},
		{"_id", storageId}})
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
	return c.JSON(http.StatusOK, storage)
}

func GetStoragesByApp(c echo.Context, client *mongo.Client, db string, storageCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	appId := c.Param("id")

	collection := client.Database(db).Collection(storageCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"user_id", userId},
		{"app_id", appId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var storageTable []utils.Storage
	err = cur.All(context.TODO(), &storageTable)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Error at getting the storages"})
	}
	return c.JSON(http.StatusOK, storageTable)
}

func UpdateStorage(c echo.Context, client *mongo.Client, db string, storageCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")
	storage := new(utils.Storage)
	if err := c.Bind(storage); err != nil {
		return err
	}
	collection := client.Database(db).Collection(storageCol)
	filter := bson.D{{"_id", storageId}, {"user_id", userId}}
	//Sanitize data
	storage.AppId = ""
	storage.UserId = ""
	storage.CreatedAt = time.Time{}
	storage.ModifiedAt = time.Now()
	storage.SharedWithId = []string{}
	storage.Type = ""
	//Create update query
	var updateFields bson.D
	tmpFields, _ := bson.Marshal(storage)
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
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Storage not found"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Storage not modified"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func DeleteStorage(c echo.Context, client *mongo.Client, db string, storageCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")
	collection := client.Database(db).Collection(storageCol)
	one, err := collection.DeleteOne(context.TODO(), bson.D{
		{"_id", storageId},
		{"user_id", userId}})
	if one.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Storage not deleted"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	//TODO: Delete the shared ids that are assigned in this storage as well from the userResourcesStratus
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func ShareStorage(c echo.Context, client *mongo.Client, db string, storageCol string, ursCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")
	body := echo.Map{}
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	if body["targetId"] == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "No Id provided"})
	}
	//TODO: Check if the target id exist as a user to share a storage
	//Set the SharedWithId in storage resource
	storageCollection := client.Database(db).Collection(storageCol)
	filter := bson.D{{"_id", storageId}, {"user_id", userId},
		{"shared_with_id", bson.D{{"$nin", bson.A{body["targetId"]}}}}}
	one, err := storageCollection.UpdateOne(context.TODO(), filter, bson.D{
		{"$push", bson.D{{"shared_with_id", body["targetId"]}}},
		{"$set", bson.D{{"modified_at", time.Now()}}}})
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Storage not found or already shared"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Storage not modified"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	//Now add the storage id into the userResourcesStatus of the target id
	ursCollection := client.Database(db).Collection(ursCol)
	filter = bson.D{{"_id", body["targetId"]},
		{"shared_storage_with_me", bson.D{{"$nin", bson.A{storageId}}}}}
	updateQuery := bson.D{{"$push", bson.D{{"shared_storage_with_me", storageId}}}}
	opts := options.Update().SetUpsert(true)
	one, err = ursCollection.UpdateOne(context.TODO(), filter, updateQuery, opts)
	if one.UpsertedCount == 0 && one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Storage already exist in urStatus"})
	}
	if err != nil {
		//TODO: Delete the target id from the storage as well cause the pipeline broke
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func UnshareStorage(c echo.Context, client *mongo.Client, db string, storageCol string, ursCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")
	body := echo.Map{}
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	if body["targetId"] == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "No Id provided"})
	}
	storageCollection := client.Database(db).Collection(storageCol)
	filter := bson.D{{"_id", storageId}, {"user_id", userId},
		{"shared_with_id", bson.D{{"$in", bson.A{body["targetId"]}}}}}
	one, err := storageCollection.UpdateOne(context.TODO(), filter, bson.D{
		{"$pull", bson.D{{"shared_with_id", body["targetId"]}}},
		{"$set", bson.D{{"modified_at", time.Now()}}}})
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Storage not found or not sharing with target"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Storage not modified"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	//Now add the storage id into the userResourcesStatus of the target id
	ursCollection := client.Database(db).Collection(ursCol)
	filter = bson.D{{"_id", body["targetId"]},
		{"shared_storage_with_me", bson.D{{"$in", bson.A{storageId}}}}}
	updateQuery := bson.D{{"$pull", bson.D{{"shared_storage_with_me", storageId}}}}
	opts := options.Update().SetUpsert(true)
	one, err = ursCollection.UpdateOne(context.TODO(), filter, updateQuery, opts)
	if one.UpsertedCount == 0 && one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Storage did not exist in urStatus"})
	}
	if err != nil {
		//TODO: Delete the target id from the storage as well cause the pipeline broke
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func AddStorageToApplication(c echo.Context, client *mongo.Client) error {
	return nil
}

func RemoveStorageFromApplication(c echo.Context, client *mongo.Client) error {
	return nil
}

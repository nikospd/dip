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
	"reflect"
	"time"
)

func CreatePullSource(c echo.Context, client *mongo.Client, db string, sourceCol string, appCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	pullSource := new(utils.PullSourceTask)
	if err := c.Bind(pullSource); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to handle request data"})
	}
	//check the provided source details
	if pullSource.SourceURI == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "No sourceURI provided"})
	}
	if pullSource.Interval < 1 || reflect.TypeOf(pullSource.Interval).String() != "int" {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Wrong interval provided"})
	}
	//check the provided appid
	if pullSource.AppId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "No application provided"})
	}
	collection := client.Database(db).Collection(appCol)
	cur := collection.FindOne(context.TODO(), bson.D{
		{"user_id", userId},
		{"_id", pullSource.AppId}})
	if cur.Err() != nil {
		if cur.Err() == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Application does not belong to the user"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create pull source"})
	}
	//TODO check if the type of the app is pull
	//Match and sanitize data
	pullSource.TaskId = utils.CreateRandomHash(20)
	pullSource.UserId = userId
	pullSource.Enabled = true
	pullSource.CreatedAt = time.Now()
	pullSource.LastExecuted = time.Time{}
	pullSource.NextExecution = time.Now()
	collection = client.Database(db).Collection(sourceCol)
	_, err := collection.InsertOne(context.TODO(), pullSource)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create pull source"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "OK", "id": pullSource.TaskId})
}

func EnablePullSource(c echo.Context, client *mongo.Client, db string, sourceCol string, status bool) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	sourceId := c.Param("id")
	collection := client.Database(db).Collection(sourceCol)
	filter := bson.D{{"_id", sourceId}, {"user_id", userId}}
	updateQuery := bson.D{{"$set", bson.D{{"enabled", status}}}}
	one, err := collection.UpdateOne(context.TODO(), filter, updateQuery)
	//Handle errors and respond
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Pull source not found"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Pull source not modified"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}
func GetPullSourceById(c echo.Context, client *mongo.Client, db string, sourceCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	sourceId := c.Param("id")
	collection := client.Database(db).Collection(sourceCol)
	one := collection.FindOne(context.TODO(), bson.D{
		{"_id", sourceId},
		{"user_id", userId}})
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad gateway"})
	}
	var source utils.PullSourceTask
	err := one.Decode(&source)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get pull source"})
	}
	return c.JSON(http.StatusOK, source)
}
func GetPullSourceByApp(c echo.Context, client *mongo.Client, db string, sourceCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	appId := c.Param("id")

	collection := client.Database(db).Collection(sourceCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"user_id", userId},
		{"app_id", appId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var sourceTable []utils.PullSourceTask
	err = cur.All(context.TODO(), &sourceTable)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Error at getting the pull sources"})
	}
	return c.JSON(http.StatusOK, sourceTable)
}
func GetPullSourceByUser(c echo.Context, client *mongo.Client, db string, sourceCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	collection := client.Database(db).Collection(sourceCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"user_id", userId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var sourceTable []utils.PullSourceTask
	err = cur.All(context.TODO(), &sourceTable)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Error at getting the pull sources"})
	}
	return c.JSON(http.StatusOK, sourceTable)
}
func UpdatePullSource(c echo.Context, client *mongo.Client, db string, sourceCol string) error {
	//update description, sourceURI and interval
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	sourceId := c.Param("id")

	source := new(utils.PullSourceTask)
	if err := c.Bind(source); err != nil {
		return err
	}
	collection := client.Database(db).Collection(sourceCol)
	filter := bson.D{{"_id", sourceId}, {"user_id", userId}}
	//Sanitize data
	source.AppId = ""
	source.UserId = ""
	source.CreatedAt = time.Time{}
	source.ModifiedAt = time.Now()
	source.LastExecuted = time.Time{}
	source.NextExecution = time.Time{}
	//Create update query and update the database
	updateQuery := bson.D{{"$set", source}}
	one, err := collection.UpdateOne(context.TODO(), filter, updateQuery)
	//Handle errors and respond
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Pull source not found"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Pull source not modified"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}
func DeletePullSource(c echo.Context, client *mongo.Client, db string, sourceCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	sourceId := c.Param("id")
	collection := client.Database(db).Collection(sourceCol)
	one, err := collection.DeleteOne(context.TODO(), bson.D{
		{"_id", sourceId},
		{"user_id", userId}})
	if one.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Pull source not deleted"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

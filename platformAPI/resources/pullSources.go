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
	//Match and sanitize data
	pullSource.TaskId = utils.CreateRandomHash(20)
	pullSource.UserId = userId
	pullSource.CreatedAt = time.Now()
	pullSource.LastExecuted = time.Time{}
	pullSource.NextExecution = time.Now()
	collection = client.Database(db).Collection(sourceCol)
	_, err := collection.InsertOne(context.TODO(), pullSource)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create pull source"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "OK"})
}

package resources

import (
	"context"
	"dev.com/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func CreateAutomation(c echo.Context, client *mongo.Client, db string, autCol string, appCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	aut := new(utils.Automation)
	if err := c.Bind(aut); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Internal server error"})
	}
	//Check necessary data
	if aut.AppId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Provide a valid app id"})
	}
	//Check if application belongs to user and update the HasIntegrations flag
	collection := client.Database(db).Collection(appCol)
	one, err := collection.UpdateOne(context.TODO(), bson.D{{"user_id", userId}, {"_id", aut.AppId}},
		bson.D{{"$set", bson.D{{"has_automations", true}}}})
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"msg": "Application does not belong to the user"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Internal server error"})
	}
	//Fill the values
	token := utils.CreateRandomHash(20)
	aut.Id = token
	aut.CreatedAt = time.Now()
	aut.UserId = userId
	//Insert document
	collection = client.Database(db).Collection(autCol)
	_, err = collection.InsertOne(context.TODO(), aut)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create resource"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "Automation created", "id": token})
}

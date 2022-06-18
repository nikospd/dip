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

func GetAutomationById(c echo.Context, client *mongo.Client, db string, autCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	automationId := c.Param("id")
	collection := client.Database(db).Collection(autCol)
	one := collection.FindOne(context.TODO(), bson.D{
		{"_id", automationId},
		{"user_id", userId}})
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad gateway"})
	}
	var aut utils.Automation
	err := one.Decode(&aut)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get automation"})
	}
	return c.JSON(http.StatusOK, aut)
}

func GetAutomationByApp(c echo.Context, client *mongo.Client, db string, autCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	appId := c.Param("id")

	collection := client.Database(db).Collection(autCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"user_id", userId},
		{"app_id", appId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var autTable []utils.Automation
	err = cur.All(context.TODO(), &autTable)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Error at getting automations"})
	}
	return c.JSON(http.StatusOK, autTable)
}

func DeleteAutomation(c echo.Context, client *mongo.Client, db string, autCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	automationId := c.Param("id")
	collection := client.Database(db).Collection(autCol)
	one, err := collection.DeleteOne(context.TODO(), bson.D{
		{"_id", automationId},
		{"user_id", userId}})
	if one.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Automation not deleted"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

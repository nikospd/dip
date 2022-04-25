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

func CreateIntegration(c echo.Context, client *mongo.Client, db string, igrCol string, appCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	igr := new(utils.Integration)
	if err := c.Bind(igr); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Internal server error"})
	}
	//Check necessary data
	if igr.AppId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Provide a valid app id"})
	}
	err := igr.CheckType()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": err.Error()})
	}
	err = igr.Option.CheckOption()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": err.Error()})
	}
	if !igr.AutomationIntegration {
		//Check if application belongs to user and update the HasIntegrations flag
		collection := client.Database(db).Collection(appCol)
		one, err := collection.UpdateOne(context.TODO(), bson.D{{"user_id", userId}, {"_id", igr.AppId}},
			bson.D{{"$set", bson.D{{"has_integrations", true}}}})
		if one.MatchedCount == 0 {
			return c.JSON(http.StatusUnauthorized, echo.Map{"msg": "Application does not belong to the user"})
		}
		if err != nil {
			return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Internal server error"})
		}
	} // Else, todo: check automations
	//Fill the values
	token := utils.CreateRandomHash(20)
	igr.Id = token
	igr.CreatedAt = time.Now()
	igr.UserId = userId
	//Insert document
	collection := client.Database(db).Collection(igrCol)
	_, err = collection.InsertOne(context.TODO(), igr)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create resource"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "Integration created", "id": token})
}

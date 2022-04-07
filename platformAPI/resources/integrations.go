package resources

import (
	"context"
	"dev.com/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func CreateIntegration(c echo.Context, client *mongo.Client, db string, igrCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	igr := new(utils.Integration)
	if err := c.Bind(igr); err != nil {
		return err
	}
	//Need check for appId if belongs to the user
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
	return c.JSON(http.StatusOK, echo.Map{"msg": "Integration created", "id": token})
}

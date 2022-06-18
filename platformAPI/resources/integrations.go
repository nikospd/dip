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

func GetIntegrationById(c echo.Context, client *mongo.Client, db string, igrCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	integrationId := c.Param("id")
	collection := client.Database(db).Collection(igrCol)
	one := collection.FindOne(context.TODO(), bson.D{
		{"_id", integrationId},
		{"user_id", userId}})
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad gateway"})
	}
	var igr utils.Integration
	err := one.Decode(&igr)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get integration"})
	}
	return c.JSON(http.StatusOK, igr)
}

func GetIntegrationByApp(c echo.Context, client *mongo.Client, db string, igrCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	appId := c.Param("id")

	collection := client.Database(db).Collection(igrCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"user_id", userId},
		{"app_id", appId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var igrTable []utils.Integration
	err = cur.All(context.TODO(), &igrTable)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Error at getting integrations"})
	}
	return c.JSON(http.StatusOK, igrTable)
}

func UpdateIntegration(c echo.Context, client *mongo.Client, db string, igrCol string) error {
	//update description and option|type
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	integrationId := c.Param("id")

	igr := new(utils.Integration)
	if err := c.Bind(igr); err != nil {
		return err
	}
	err := igr.CheckType()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": err.Error()})
	}
	err = igr.Option.CheckOption()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": err.Error()})
	}
	collection := client.Database(db).Collection(igrCol)
	filter := bson.D{{"_id", integrationId}, {"user_id", userId}}
	//Sanitize data
	igr.AppId = ""
	igr.UserId = ""
	igr.CreatedAt = time.Time{}
	igr.ModifiedAt = time.Now()
	//Create update query and update the database
	updateQuery := bson.D{{"$set", igr}}
	one, err := collection.UpdateOne(context.TODO(), filter, updateQuery)
	//Handle errors and respond
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Integration not found"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Integration not modified"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func DeleteIntegration(c echo.Context, client *mongo.Client, db string, igrCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	integrationId := c.Param("id")
	collection := client.Database(db).Collection(igrCol)
	one, err := collection.DeleteOne(context.TODO(), bson.D{
		{"_id", integrationId},
		{"user_id", userId}})
	if one.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Integration not deleted"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

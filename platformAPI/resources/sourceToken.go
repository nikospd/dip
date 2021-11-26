package resources

import (
	"context"
	"dev.com/utils"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// TODO: check for missing variables

func CreateSourceToken(c echo.Context, client *mongo.Client, db string, tokenCol string, appCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	//Create token
	stc := new(utils.SourceTokenClaims)
	if err := c.Bind(stc); err != nil {
		return err
	}
	stc.UserId = userId
	stc.CreatedAt = time.Now()
	token := utils.CreateRandomHash(20)
	if token == "" {
		return errors.New("error on creation of hash source token")
	}
	stc.SourceToken = token
	//Search if application belongs to this user
	collection := client.Database(db).Collection(appCol)
	cur := collection.FindOne(context.TODO(), bson.D{
		{"user_id", userId},
		{"_id", stc.AppId}})
	if cur.Err() != nil {
		if cur.Err() == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Application does not belong to the user"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create token"})
	}
	//TODO check if app is for push source
	//Add token to the database
	collection = client.Database(db).Collection(tokenCol)
	_, err := collection.InsertOne(context.TODO(), stc)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create token"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "OK", "id": stc.SourceToken})
}

func GetSourceTokenByUser(c echo.Context, client *mongo.Client, db string, tokenCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id

	collection := client.Database(db).Collection(tokenCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"user_id", userId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var stcTable []utils.SourceTokenClaims
	cur.All(context.TODO(), &stcTable)
	return c.JSON(http.StatusOK, stcTable)
}

func GetSourceTokenById(c echo.Context, client *mongo.Client, db string, tokenCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	sourceId := c.Param("id")
	collection := client.Database(db).Collection(tokenCol)
	one := collection.FindOne(context.TODO(), bson.D{
		{"_id", sourceId}, {"user_id", userId}})
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var stc utils.SourceTokenClaims
	one.Decode(&stc)
	return c.JSON(http.StatusOK, stc)
}

func GetSourceTokenByApp(c echo.Context, client *mongo.Client, db string, tokenCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	appId := c.Param("id")
	collection := client.Database(db).Collection(tokenCol)
	cur, err := collection.Find(context.TODO(), bson.D{
		{"user_id", userId},
		{"app_id", appId}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("not a match!")
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var stcTable []utils.SourceTokenClaims
	cur.All(context.TODO(), &stcTable)
	return c.JSON(http.StatusOK, stcTable)
}

func ModifySourceToken(c echo.Context, client *mongo.Client, db string, tokenCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	tokenId := c.Param("id")
	stc := new(utils.SourceTokenClaims)
	if err := c.Bind(stc); err != nil {
		return err
	}
	collection := client.Database(db).Collection(tokenCol)
	filter := bson.D{{"_id", tokenId}, {"user_id", userId}}
	//Sanitize data
	stc.SourceToken = ""
	stc.AppId = ""
	stc.UserId = ""
	stc.CreatedAt = time.Time{}
	stc.ModifiedAt = time.Now()
	//Create update query
	var updateFields bson.D
	tmpFields, _ := bson.Marshal(stc)
	unmarshalErr := bson.Unmarshal(tmpFields, &updateFields)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	updateQuery := bson.D{{"$set", updateFields}}
	one, err := collection.UpdateOne(context.TODO(), filter, updateQuery)
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Token not found"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Token not modified"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func DeleteSourceToken(c echo.Context, client *mongo.Client, db string, tokenCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	tokenId := c.Param("id")
	collection := client.Database(db).Collection(tokenCol)
	one, err := collection.DeleteOne(context.TODO(), bson.D{{"_id", tokenId}, {"user_id", userId}})
	if one.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Token not deleted"})
	}
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

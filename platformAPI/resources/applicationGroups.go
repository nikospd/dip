package resources

import (
	"context"
	"dev.com/utils"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func CreateApplicationGroup(c echo.Context, client *mongo.Client, db string, groupCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	//Create group
	group := new(utils.ApplicationGroup)
	if err := c.Bind(group); err != nil {
		return err
	}
	group.UserId = userId
	group.CreatedAt = time.Now()
	token := utils.CreateRandomHash(20)
	if token == "" {
		return errors.New("error on creation of hash application group")
	}
	group.GroupId = token
	//Sanitize data
	group.Applications = []string{}
	group.ModifiedAt = time.Time{}
	collection := client.Database(db).Collection(groupCol)
	_, err := collection.InsertOne(context.TODO(), group)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create resource"})
	}
	return c.JSON(http.StatusCreated, echo.Map{"msg": "Application group created", "id": group.GroupId})
}

func GetApplicationGroupById(c echo.Context, client *mongo.Client, db string, groupCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	groupId := c.Param("id")
	collection := client.Database(db).Collection(groupCol)
	one := collection.FindOne(context.TODO(), bson.D{{"user_id", userId}, {"_id", groupId}})
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad gateway"})
	}
	var group utils.ApplicationGroup
	err := one.Decode(&group)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get application group"})
	}
	if group.Applications != nil {
		group.NumOfApplications = len(group.Applications)
	}
	return c.JSON(http.StatusOK, group)
}

func GetApplicationGroupByUser(c echo.Context, client *mongo.Client, db string, groupCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id

	findQuery := bson.D{{"user_id", userId}}
	collection := client.Database(db).Collection(groupCol)
	cur, err := collection.Find(context.TODO(), findQuery)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var groupTable []utils.ApplicationGroup
	err = cur.All(context.TODO(), &groupTable)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Error at getting the application groups"})
	}
	for idx, group := range groupTable {
		if group.Applications != nil {
			groupTable[idx].NumOfApplications = len(group.Applications)
		}
	}
	return c.JSON(http.StatusOK, groupTable)
}

func ChangeApplicationGroup(c echo.Context, client *mongo.Client, db string, groupCol string, appCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	newGroupId := c.Param("id")
	appId := c.QueryParam("appId")
	if appId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Please provide a valid application id"})
	}
	if newGroupId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Please provide a valid application group id as a target"})
	}

	collection := client.Database(db).Collection(appCol)
	oneFind := collection.FindOne(context.TODO(), bson.D{{"_id", appId}, {"user_id", userId}})
	var app utils.Application
	if oneFind.Err() != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to fetch application"})
	}
	err := oneFind.Decode(&app)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to fetch application"})
	}
	if newGroupId == app.ApplicationGroupId {
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "New group id cannot be the same with the old one"})
	}
	collection = client.Database(db).Collection(groupCol)
	oneUpdate, err := collection.UpdateOne(context.TODO(),
		bson.D{{"user_id", userId}, {"_id", newGroupId}},
		bson.D{{"$push", bson.D{{"applications", appId}}}})
	if oneUpdate.ModifiedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Target app group does not belongs to user"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to add application to new group"})
	}
	oneUpdate, err = collection.UpdateOne(context.TODO(),
		bson.D{{"user_id", userId}, {"_id", app.ApplicationGroupId}},
		bson.D{{"$pull", bson.D{{"applications", appId}}}})
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to remove application from old group"})
		//	todo: add strategy to be fault tolerance
	}
	collection = client.Database(db).Collection(appCol)
	_, err = collection.UpdateOne(context.TODO(), bson.D{{"_id", appId}},
		bson.D{{"$set", bson.D{{"application_group_id", newGroupId}}}})
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to change application to new group"})
		//	todo: add strategy to be fault tolerance
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func DeleteApplicationGroup(c echo.Context, client *mongo.Client, db string, groupCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	groupId := c.Param("id")

	collection := client.Database(db).Collection(groupCol)
	oneFind := collection.FindOne(context.TODO(), bson.D{{"_id", groupId}, {"user_id", userId}})
	if oneFind.Err() != nil {
		if oneFind.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad gateway"})
	}
	var group utils.ApplicationGroup
	err := oneFind.Decode(&group)
	if len(group.Applications) != 0 {
		return c.JSON(http.StatusNotAcceptable, echo.Map{"msg": "Application group is not empty"})
	}
	one, err := collection.DeleteOne(context.TODO(), bson.D{
		{"_id", groupId},
		{"user_id", userId}})
	if one.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Group not deleted"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

func UpdateApplicationGroup(c echo.Context, client *mongo.Client, db string, groupCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	groupId := c.Param("id")
	group := new(utils.ApplicationGroup)
	if err := c.Bind(group); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	collection := client.Database(db).Collection(groupCol)
	filter := bson.D{{"_id", groupId}, {"user_id", userId}}
	//Sanitize data
	group.GroupId = groupId
	group.UserId = ""
	group.CreatedAt = time.Time{}
	group.ModifiedAt = time.Now()
	group.Applications = []string{}
	//Update
	one, err := collection.UpdateOne(context.TODO(), filter, bson.D{{"$set", group}})
	//Handle errors and respond
	if one.MatchedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"msg": "Application group not found"})
	}
	if one.ModifiedCount == 0 {
		return c.JSON(http.StatusNotModified, echo.Map{"msg": "Application group not modified"})
	}
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad Gateway"})
	}
	return c.JSON(http.StatusOK, echo.Map{"msg": "OK"})
}

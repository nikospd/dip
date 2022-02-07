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

// todo: remove application from group on delete
// todo: Change group of one application

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

//func addApplicationToGroup(c echo.Context, client *mongo.Client, db string, groupCol string, appCol string) error {
//	user := c.Get("user").(*jwt.Token)
//	claims := user.Claims.(*jwt.StandardClaims)
//	userId := claims.Id
//	groupId := c.Param("id")
//	appId := c.QueryParam("appId")
//
//}

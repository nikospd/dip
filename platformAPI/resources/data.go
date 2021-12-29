package resources

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
)

func GetStorageData(c echo.Context, client *mongo.Client, resourcesDb string, dataDb string, storageCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	userId := claims.Id
	storageId := c.Param("id")
	//Var for paginated data
	page := 1
	nPerPage := 10 //Documents per page
	var err error
	if c.QueryParam("page") != "" {
		page, err = strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Wrong page given"})
		}
	}
	if c.QueryParam("size") != "" {
		nPerPage, err = strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Wrong size given"})
		}
	}
	//Get storage from id
	//Check if storage belongs to this userId or if it is shared to it
	findQuery := bson.D{
		{"_id", storageId},
		{"$or", bson.A{
			bson.D{{"user_id", userId}},
			bson.D{{"shared_with_id", bson.D{{"$in", bson.A{userId}}}}}},
		}}
	storageCollection := client.Database(resourcesDb).Collection(storageCol)
	one := storageCollection.FindOne(context.TODO(), findQuery)
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found Or storage does not belong to the user"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Bad gateway"})
	}
	//Get data from that storage
	dataCollection := client.Database(dataDb).Collection(storageId)
	opts := options.Find()
	opts.Projection = bson.D{{"_id", 0}}
	opts.SetLimit(int64(nPerPage))
	opts.SetSkip(int64(nPerPage * (page - 1)))
	cur, dferr := dataCollection.Find(context.TODO(), bson.D{}, opts)
	if dferr != nil {
		if dferr == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var dataTable []map[string]interface{}
	cur.All(context.TODO(), &dataTable)
	numDoc, _ := dataCollection.CountDocuments(context.TODO(), bson.D{})
	return c.JSON(http.StatusOK, echo.Map{"data": dataTable, "totalDocs": numDoc})
}

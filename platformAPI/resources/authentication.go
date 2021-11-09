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
)

func UserLogin(c echo.Context, client *mongo.Client, signingKey []byte, db string, col string) error {
	//TODO: Use encryption for passwords
	collection := client.Database(db).Collection(col)
	credentials := new(utils.LoginUserCredentials)
	if err := c.Bind(credentials); err != nil {
		return err
	}
	cur := collection.FindOne(context.TODO(), bson.D{
		{"username", credentials.Username},
		{"password", credentials.Password}})
	if cur.Err() != nil {
		if cur.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var user utils.LoginUserCredentials
	err := cur.Decode(&user)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to handle credentials"})
	}
	// Create the Claims
	claims := &jwt.StandardClaims{
		Id: user.UserId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"token": t})
}

func UserRegister() {
	fmt.Println("Implement me")
}

func UserForgotPassword() {
	fmt.Println("Implement me")
}

func UserChangePassword() {
	fmt.Println("Implement me")
}

func GetUser(c echo.Context, client *mongo.Client, db string, userCol string) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.StandardClaims)
	searchValue := claims.Id //userid
	searchKey := "_id"
	username := c.QueryParam("name")
	if username != "" {
		searchValue = username
		searchKey = "username"
	}
	userCollection := client.Database(db).Collection(userCol)
	one := userCollection.FindOne(context.TODO(), bson.D{{searchKey, searchValue}})
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var profile utils.UserProfile
	err := one.Decode(&profile)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get profile"})
	}
	return c.JSON(http.StatusOK, profile)
}

func GetUserByMail() {
	fmt.Println("Implement me")
}

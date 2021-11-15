package resources

import (
	"context"
	"crypto/md5"
	"dev.com/utils"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func UserLogin(c echo.Context, client *mongo.Client, signingKey []byte, db string, col string) error {
	collection := client.Database(db).Collection(col)
	credentials := new(utils.User)
	if err := c.Bind(credentials); err != nil {
		return err
	}
	hash := md5.Sum([]byte(credentials.Password))
	credentials.Password = hex.EncodeToString(hash[:])
	cur := collection.FindOne(context.TODO(), bson.D{
		{"username", credentials.Username},
		{"password", credentials.Password}})
	if cur.Err() != nil {
		if cur.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var user utils.User
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

func UserRegister(c echo.Context, col *mongo.Collection) error{
	user := new(utils.User)
	err := c.Bind(user)
	if err != nil{
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get the user at register"})
	}
	//Check the input params
	if user.Username == ""{
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Username not provided"})
	}
	if user.Password == ""{
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Password not provided"})
	}
	if user.Email == ""{
		return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Email not provided"})
	}
	user.UserId = utils.CreateRandomHash(20)
	hash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])
	user.CreatedAt = time.Now()
	_, err = col.InsertOne(context.TODO(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.JSON(http.StatusBadRequest, echo.Map{"msg": "Username or email already exist"})
		}
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to create user"})
	}
	return c.JSON(http.StatusOK, echo.Map{"userId": user.UserId, "msg": "OK"})
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
	email := c.QueryParam("email")
	if username != "" {
		searchValue = username
		searchKey = "username"
	} else if email!=""{
		searchValue = email
		searchKey = "email"
	}
	userCollection := client.Database(db).Collection(userCol)
	opt:=options.FindOne()
	opt.Projection = bson.D{{"password", 0}}
	one := userCollection.FindOne(context.TODO(),
		bson.D{{searchKey, searchValue}}, opt)
	if one.Err() != nil {
		if one.Err() == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, echo.Map{"msg": "Not Found"})
		}
	}
	var profile utils.User
	err := one.Decode(&profile)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"msg": "Failed to get profile"})
	}
	return c.JSON(http.StatusOK, profile)
}
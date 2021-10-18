package main

import (
	"context"
	"dev.com/resources"
	"dev.com/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//TODO: fix loggers
//TODO: add configuration files
//TODO: error handling

var client *mongo.Client
var mySigningKey []byte
var cfg utils.Configuration

func main() {
	/*
		Read configuration file
	*/
	utils.ReadConf("config.json", &cfg)
	/*
		Connect to MongoDB
	*/
	mongoCred := cfg.MongoCredentials
	mongoUri := utils.MongoCredentials(mongoCred.User, mongoCred.Password, mongoCred.Host, mongoCred.Port)
	clientOptions := options.Client().ApplyURI(mongoUri)
	var connectionError error
	client, connectionError = mongo.Connect(context.TODO(), clientOptions)
	if connectionError != nil {
		log.Fatalln(connectionError)
	}
	connectionError = client.Ping(context.TODO(), nil)
	if connectionError != nil {
		log.Fatalln(connectionError)
	}
	/*
		Prepare the web server
	*/
	e := echo.New()
	/*
		Middleware
	*/
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	/*
		Assign resources to restricted endpoints
	*/
	mySigningKey = []byte(cfg.SigningKey)
	jwtConfig := middleware.JWTConfig{
		Claims:     &jwt.StandardClaims{},
		SigningKey: mySigningKey,
	}
	r := e.Group("/")
	r.Use(middleware.JWTWithConfig(jwtConfig))
	//SourceToken routing
	r.POST("source/token", createSourceToken)
	r.GET("source/token", getSourceTokenByUser)
	r.GET("source/token/:id", getSourceTokenByApp)
	r.PUT("source/token/:id", modifySourceToken)
	r.DELETE("source/token/:id", deleteSourceToken)
	//Application property routing
	r.POST("application", createApplication)
	r.GET("application", getApplicationsByUser)
	r.PUT("application/:id", updateApplication)
	r.DELETE("application/:id", deleteApplication)
	//Storage property routing
	r.POST("storage", createStorage)
	/*
		Assign resources to unauthenticated endpoints
	*/
	e.POST("/login", userLogin)
	/*
		Start the web server
	*/
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.ApiPort)))
}

/*
Wrapper functions for resources
*/
func userLogin(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	userCol := cfg.MongoCollection.Users
	return resources.UserLogin(c, client, mySigningKey, db, userCol)
}
func createSourceToken(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	tokenCol := cfg.MongoCollection.SourceTokens
	appCol := cfg.MongoCollection.Applications
	return resources.CreateSourceToken(c, client, db, tokenCol, appCol)
}
func getSourceTokenByUser(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	tokenCol := cfg.MongoCollection.SourceTokens
	return resources.GetSourceTokenByUser(c, client, db, tokenCol)
}
func getSourceTokenByApp(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	tokenCol := cfg.MongoCollection.SourceTokens
	return resources.GetSourceTokenByApp(c, client, db, tokenCol)
}
func modifySourceToken(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	tokenCol := cfg.MongoCollection.SourceTokens
	return resources.ModifySourceToken(c, client, db, tokenCol)
}
func deleteSourceToken(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	tokenCol := cfg.MongoCollection.SourceTokens
	return resources.DeleteSourceToken(c, client, db, tokenCol)
}
func createApplication(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	appCol := cfg.MongoCollection.Applications
	return resources.CreateApplication(c, client, db, storageCol, appCol)
}
func getApplicationsByUser(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	appCol := cfg.MongoCollection.Applications
	return resources.GetApplicationsByUser(c, client, db, appCol)
}
func updateApplication(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	appCol := cfg.MongoCollection.Applications
	return resources.UpdateApplication(c, client, db, appCol)
}
func deleteApplication(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	appCol := cfg.MongoCollection.Applications
	return resources.DeleteApplication(c, client, db, appCol)
}
func createStorage(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	return resources.CreateStorage(c, client, db, storageCol)
}

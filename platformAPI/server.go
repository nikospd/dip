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
		Swagger
	*/
	e.Static("/swagger-ui.css", "swaggerui/swagger-ui.css")
	e.Static("/swagger-ui-bundle.js", "swaggerui/swagger-ui-bundle.js")
	e.Static("/swagger-ui-standalone-preset.js", "swaggerui/swagger-ui-standalone-preset.js")
	e.Static("/swagger.json", "swaggerui/swagger.json")
	e.Static("/api", "swaggerui/index.html")
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
	//User Info routing
	r.GET("user/profile", getUser)
	//SourceToken routing
	r.POST("source/token", createSourceToken)
	r.GET("source/tokens", getSourceTokenByUser)
	r.GET("source/token/:id", getSourceTokenById)
	r.GET("source/tokens/:id", getSourceTokenByApp)
	r.PUT("source/token/:id", modifySourceToken)
	r.DELETE("source/token/:id", deleteSourceToken)
	//Pull source routing
	//TODO: enable/disable/get/put/delete
	r.POST("source/pull", createPullSource)
	//Application property routing
	r.POST("application", createApplication)
	r.GET("application", getApplicationsByUser)
	r.GET("application/:id", getApplicationById)
	r.PUT("application/:id", updateApplication)
	r.DELETE("application/:id", deleteApplication)
	//Storage property routing
	r.POST("storage", createStorage)
	r.GET("storage/:id", getStorageById)
	r.GET("storages/:id", getStorageByApp)
	r.GET("storages", getStoragesByUser)
	r.PUT("storage/:id", updateStorage)
	r.DELETE("storage/:id", deleteStorage)
	r.POST("storage/share/:id", shareStorage)
	r.POST("storage/unshare/:id", unshareStorage)
	r.POST("storage/attach/:id", attachStorage)
	r.POST("storage/detach/:id", detachStorage)
	//Data routing
	r.GET("storage/data/:id", getStorageData)
	/*
		Assign resources to unauthenticated endpoints
	*/
	e.POST("/user/login", userLogin)
	e.POST("/user/register", userRegister)
	/*
		Handle cors policy
		Normally you should add cors with config and add your domain here.
	*/
	e.Use(middleware.CORS())
	/*
		Start the web server
	*/
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.ApiPort)))
}

/*
Wrapper functions for resources
*/
func getUser(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	userCol := cfg.MongoCollection.Users
	return resources.GetUser(c, client, db, userCol)
}
func userLogin(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	userCol := cfg.MongoCollection.Users
	return resources.UserLogin(c, client, mySigningKey, db, userCol)
}
func userRegister(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	userCol := cfg.MongoCollection.Users
	col := client.Database(db).Collection(userCol)
	return resources.UserRegister(c, col)
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
func getSourceTokenById(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	tokenCol := cfg.MongoCollection.SourceTokens
	return resources.GetSourceTokenById(c, client, db, tokenCol)
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
func createPullSource(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	sourceCol := cfg.MongoCollection.PullSources
	appCol := cfg.MongoCollection.Applications
	return resources.CreatePullSource(c, client, db, sourceCol, appCol)
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
func getApplicationById(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	appCol := cfg.MongoCollection.Applications
	return resources.GetApplicationById(c, client, db, appCol)
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
func getStorageByApp(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	return resources.GetStoragesByApp(c, client, db, storageCol)
}
func getStoragesByUser(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	return resources.GetStoragesByUser(c, client, db, storageCol)
}
func getStorageById(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	return resources.GetStorageById(c, client, db, storageCol)
}
func updateStorage(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	return resources.UpdateStorage(c, client, db, storageCol)
}
func deleteStorage(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	return resources.DeleteStorage(c, client, db, storageCol)
}
func shareStorage(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	ursCol := cfg.MongoCollection.URStatus
	return resources.ShareStorage(c, client, db, storageCol, ursCol)
}
func unshareStorage(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	ursCol := cfg.MongoCollection.URStatus
	return resources.UnshareStorage(c, client, db, storageCol, ursCol)
}
func attachStorage(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	appCol := cfg.MongoCollection.Applications
	return resources.AttachStorage(c, client, db, storageCol, appCol)
}
func detachStorage(c echo.Context) error {
	db := cfg.MongoDatabase.Resources
	storageCol := cfg.MongoCollection.Storages
	appCol := cfg.MongoCollection.Applications
	return resources.DetachStorage(c, client, db, storageCol, appCol)
}
func getStorageData(c echo.Context) error {
	resourcesDb := cfg.MongoDatabase.Resources
	dataDb := cfg.MongoDatabase.Data
	storageCol := cfg.MongoCollection.Storages
	return resources.GetStorageData(c, client, resourcesDb, dataDb, storageCol)
}
